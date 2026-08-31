package gate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Clock is the seam that lets tests advance time without sleeping (§22) — the
// only two operations the gate ever needs from a clock.
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

// SystemClock is the production Clock: the real wall clock.
type SystemClock struct{}

func (SystemClock) Now() time.Time                         { return time.Now() }
func (SystemClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Observation is one successful poll of the state endpoint for a single rule.
type Observation struct {
	Rules      []StateRule   // may be empty — an authoritative 2xx saying the rule is absent (§14.5)
	GrafanaNow time.Time     // the Date header — H4
	Skew       time.Duration // serverDate - (t_send+t_headers)/2, signed (§16)
	SkewBound  time.Duration // (t_headers-t_send)/2 — RTT/2 to the response headers
	Latency    time.Duration // t_send through the full body read — see requestResult.Latency
}

// TransportError marks a failure worth retrying: a non-2xx response, a
// network failure, or a body that failed to parse. It is never a deleted rule
// (an authoritative 2xx with no matching rule is not this) and never a clock
// problem (a missing/unparseable Date header or an out-of-bounds skew is a
// hard error instead — see doRequest). Never conflate them (§14.5).
type TransportError struct {
	Err    error
	Status int // 0 when the failure never got a status (network/transport failure)
}

func (e *TransportError) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("transport error: status %d: %v", e.Status, e.Err)
	}
	return fmt.Sprintf("transport error: %v", e.Err)
}

func (e *TransportError) Unwrap() error { return e.Err }

// RetryExhaustedError is what retryTransport returns once it gives up after
// too many sequential *TransportError failures. It deliberately does not
// implement Unwrap into the underlying *TransportError: once retries are
// exhausted the result is a hard, terminal failure, and
// errors.AsType[*TransportError] must never re-classify it as retryable —
// that is the exact conflation §19.3 case 1 forbids. Cause is still exposed
// as a plain field (and folded into Error()'s text) so a caller can log or
// inspect it; it just cannot flow back into the retry classification.
type RetryExhaustedError struct {
	Failures int
	Cause    error
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("gave up after %d sequential failures: %v", e.Failures, e.Cause)
}

// Source is everything the gate reads from Grafana. httpSource is the one
// production implementation; every later phase's tests use a scripted fake
// (source_fake_test.go) instead of real HTTP.
type Source interface {
	Version(ctx context.Context) (string, error)
	Definitions(ctx context.Context) ([]Definition, error)
	RuleState(ctx context.Context, title string) (Observation, error)
}

// grafanaVersion is a parsed major.minor.patch triple.
type grafanaVersion struct{ major, minor, patch int }

func (v grafanaVersion) String() string { return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch) }

// ord encodes the triple as a single comparable integer. Safe as long as
// minor and patch stay under 1000, true of every real Grafana version.
func (v grafanaVersion) ord() int64 {
	return int64(v.major)*1_000_000 + int64(v.minor)*1_000 + int64(v.patch)
}

func parseGrafanaVersion(s string) (grafanaVersion, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return grafanaVersion{}, errors.New("empty version string")
	}
	// Grafana's /api/health always reports exactly three components
	// (health.json: "13.1.0"). Require all three explicitly rather than
	// defaulting missing ones to zero or silently dropping extras — either
	// would accept a value ("13", "13.1.0.5") that was never actually seen
	// and never verified against.
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return grafanaVersion{}, fmt.Errorf("unparseable version %q: want exactly 3 dot-separated components, got %d", s, len(parts))
	}
	var v grafanaVersion
	fields := [3]*int{&v.major, &v.minor, &v.patch}
	for i, field := range fields {
		// Trim any trailing non-digit suffix (prerelease/build metadata, e.g.
		// "0+security") rather than requiring an exact numeric match.
		digits := parts[i]
		j := 0
		for j < len(digits) && digits[j] >= '0' && digits[j] <= '9' {
			j++
		}
		if j == 0 {
			return grafanaVersion{}, fmt.Errorf("unparseable version %q", s)
		}
		n, err := strconv.Atoi(digits[:j])
		if err != nil {
			return grafanaVersion{}, fmt.Errorf("unparseable version %q: %w", s, err)
		}
		*field = n
	}
	return v, nil
}

// supportedGrafanaMin and supportedGrafanaMax bound the platform this gate is
// verified against (§2.7 control 2, §21.5): >= 13.0.0, < 14.0.0.
var (
	supportedGrafanaMin = grafanaVersion{13, 0, 0}
	supportedGrafanaMax = grafanaVersion{14, 0, 0} // exclusive
)

// CheckGrafanaVersion enforces the supported range. An unparseable or
// out-of-range version is a hard error naming both what was found and what is
// supported — trusting an unverified schema is exactly the deprecation risk
// §2.7 control 2 exists to catch.
func CheckGrafanaVersion(version string) error {
	v, err := parseGrafanaVersion(version)
	if err != nil {
		return fmt.Errorf("grafana version %q: %w (supported: >=%s, <%s)",
			version, err, supportedGrafanaMin, supportedGrafanaMax)
	}
	if v.ord() < supportedGrafanaMin.ord() || v.ord() >= supportedGrafanaMax.ord() {
		return fmt.Errorf("unsupported grafana version %q (supported: >=%s, <%s)",
			version, supportedGrafanaMin, supportedGrafanaMax)
	}
	return nil
}

// httpSource is the production Source: stdlib net/http only, bearer auth
// from a token supplied at construction (the caller reads it from the
// environment — §20.2 — this type never touches env itself), and manual
// strict decoding via ParseState/ParseDefinitions (H1). The retry limit and
// backoff parameters are struct fields with production defaults set here,
// not package constants, so a test can shrink them without a hook.
type httpSource struct {
	baseURL string
	token   string
	client  *http.Client
	clock   Clock

	maxSequentialFailures int
	backoffBase           time.Duration
	backoffCap            time.Duration
}

// NewHTTPSource builds the production Source. token is never logged and
// never enters an error string (§20.2) — it is used only to set the
// Authorization header.
func NewHTTPSource(baseURL, token string, clock Clock) Source {
	return &httpSource{
		baseURL:               strings.TrimSuffix(baseURL, "/"),
		token:                 token,
		client:                &http.Client{Timeout: 30 * time.Second},
		clock:                 clock,
		maxSequentialFailures: 5,
		backoffBase:           time.Second,
		backoffCap:            30 * time.Second,
	}
}

func (s *httpSource) Version(ctx context.Context) (string, error) {
	return retryTransport(ctx, s.clock, s.maxSequentialFailures, s.backoffBase, s.backoffCap, func() (string, error) {
		r, err := s.doRequest(ctx, "/api/health")
		if err != nil {
			return "", err
		}
		var health struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(r.Body, &health); err != nil {
			return "", &TransportError{Err: fmt.Errorf("parse /api/health: %w", err)}
		}
		if health.Version == "" {
			return "", &TransportError{Err: errors.New("/api/health: empty version")}
		}
		return health.Version, nil
	})
}

func (s *httpSource) Definitions(ctx context.Context) ([]Definition, error) {
	return retryTransport(ctx, s.clock, s.maxSequentialFailures, s.backoffBase, s.backoffCap, func() ([]Definition, error) {
		r, err := s.doRequest(ctx, "/api/ruler/grafana/api/v1/rules")
		if err != nil {
			return nil, err
		}
		defs, parseErr := ParseDefinitions(r.Body)
		if parseErr != nil {
			return nil, &TransportError{Err: fmt.Errorf("parse ruler definitions: %w", parseErr)}
		}
		return defs, nil
	})
}

func (s *httpSource) RuleState(ctx context.Context, title string) (Observation, error) {
	path := "/api/prometheus/grafana/api/v1/rules?rule_name=" + url.QueryEscape(title)
	return retryTransport(ctx, s.clock, s.maxSequentialFailures, s.backoffBase, s.backoffCap, func() (Observation, error) {
		r, err := s.doRequest(ctx, path)
		if err != nil {
			return Observation{}, err
		}
		rules, parseErr := ParseState(r.Body)
		if parseErr != nil {
			// Treated as transient, not a schema break: an unparseable 2xx
			// is far more likely a mid-stream hiccup than a permanent shape
			// change, and H1's strict parser already turns a real shape
			// change into a loud per-field error the moment it's visible.
			return Observation{}, &TransportError{Err: fmt.Errorf("parse rule state: %w", parseErr)}
		}
		return Observation{
			Rules:      rules,
			GrafanaNow: r.ServerDate,
			Skew:       r.Skew,
			SkewBound:  r.SkewBound,
			Latency:    r.Latency,
		}, nil
	})
}

// requestResult is the outcome of one successful HTTP attempt in doRequest:
// the raw body plus everything derived from timing the round trip against
// the response's own clock (§16).
type requestResult struct {
	Body       []byte
	ServerDate time.Time     // the Date header — H4
	Skew       time.Duration // serverDate - (t_send+t_headers)/2, signed
	SkewBound  time.Duration // (t_headers-t_send)/2 — RTT/2 to the response headers
	// Latency spans t_send through the full body read (§5.2's budget check
	// needs the whole poll's wall time, or a schedule feasibility check that
	// only sees header latency goes optimistic — fail-open). It does not
	// include the caller's subsequent JSON parse (ParseState/ParseDefinitions
	// run outside doRequest); if P4's budget accounting needs parse time
	// folded in too, extend here rather than approximating it at the call
	// site.
	Latency time.Duration
}

// doRequest performs one HTTP GET and classifies the outcome (§14.5, §16):
// a network failure, a non-2xx status, or a body-read failure is retryable
// (*TransportError); a missing or unparseable Date header, or a skew beyond
// SkewHardLimit, is a hard error — retrying can never fix either, so neither
// may enter the backoff loop (H4).
//
// The Date-header/skew check runs for every endpoint this hits, including
// /api/health — broader than §16's own scope, which only discusses the state
// endpoint. Deliberate: a skewed clock discovered only once RuleState starts
// polling is a skew that has already masked whatever /api/health and the
// ruler read reported; failing closed at the first response catches it
// before any of that is trusted, and every response comes with a Date header
// for free.
func (s *httpSource) doRequest(ctx context.Context, path string) (requestResult, error) {
	req, buildErr := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+path, nil)
	if buildErr != nil {
		return requestResult{}, fmt.Errorf("build request for %s: %w", path, buildErr)
	}
	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	}

	tSend := s.clock.Now()
	resp, doErr := s.client.Do(req)
	tHeaders := s.clock.Now()
	if doErr != nil {
		return requestResult{}, &TransportError{Err: doErr}
	}
	defer resp.Body.Close()

	b, readErr := io.ReadAll(resp.Body)
	tBodyRead := s.clock.Now()
	if readErr != nil {
		return requestResult{}, &TransportError{Err: fmt.Errorf("read response body (status %d): %w", resp.StatusCode, readErr)}
	}
	latency := tBodyRead.Sub(tSend)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return requestResult{}, &TransportError{Err: fmt.Errorf("unexpected status %d", resp.StatusCode), Status: resp.StatusCode}
	}

	dateHeader := resp.Header.Get("Date")
	if dateHeader == "" {
		return requestResult{}, fmt.Errorf("%s: response has no Date header (H4)", path)
	}
	serverDate, parseErr := http.ParseTime(dateHeader)
	if parseErr != nil {
		return requestResult{}, fmt.Errorf("%s: unparseable Date header %q: %w", path, dateHeader, parseErr)
	}

	bound := tHeaders.Sub(tSend) / 2
	mid := tSend.Add(bound)
	signedSkew := serverDate.Sub(mid)
	absSkew := signedSkew
	if absSkew < 0 {
		absSkew = -absSkew
	}
	if absSkew > SkewHardLimit {
		return requestResult{}, fmt.Errorf("%s: clock skew %s exceeds hard limit %s (§16)", path, absSkew, SkewHardLimit)
	}

	return requestResult{Body: b, ServerDate: serverDate, Skew: signedSkew, SkewBound: bound, Latency: latency}, nil
}

// retryTransport runs fn, retrying with backoff only while it fails with a
// *TransportError — any other error is a hard error and returns immediately,
// never retried. failures counts consecutive *TransportError results;
// exceeding maxFailures gives up with a wrapped hard error (§19.3 case 1).
// The wait between attempts goes through clock.After so a test with a fake
// Clock never sleeps on real time (§22).
func retryTransport[T any](ctx context.Context, clock Clock, maxFailures int, backoffBase, backoffCap time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	failures := 0
	for {
		v, err := fn()
		if err == nil {
			return v, nil
		}
		if _, ok := errors.AsType[*TransportError](err); !ok {
			return zero, err
		}
		failures++
		if failures > maxFailures {
			return zero, &RetryExhaustedError{Failures: failures, Cause: err}
		}
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-clock.After(backoffDelay(backoffBase, backoffCap, failures)):
		}
	}
}

// backoffDelay is 1s base, doubling per failure, capped at maxDelay, with
// ±20% jitter (§5's filled-in value for maxSequentialFailures).
func backoffDelay(base, maxDelay time.Duration, failureCount int) time.Duration {
	d := base
	for i := 1; i < failureCount && d < maxDelay; i++ {
		d *= 2
	}
	if d > maxDelay {
		d = maxDelay
	}
	jitter := 0.8 + rand.Float64()*0.4 // [0.8, 1.2]
	return time.Duration(float64(d) * jitter)
}
