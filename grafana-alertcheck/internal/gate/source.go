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

// Clock is the seam that lets tests advance time without sleeping — the only
// two operations the gate ever needs from a clock.
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
	Rules      []StateRule   // may be empty — an authoritative 2xx saying the rule is absent
	GrafanaNow time.Time     // the response's Date header
	Skew       time.Duration // serverDate - (t_send+t_headers)/2, signed
	SkewBound  time.Duration // (t_headers-t_send)/2 — RTT/2 to the response headers
	Latency    time.Duration // t_send through the full body read — see requestResult.Latency
}

// TransportError marks a failure worth retrying: a non-2xx response, a network
// failure, or a body that failed to parse. Not a deleted rule (an authoritative
// 2xx) and not a clock problem (a hard error — see doRequest).
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

// RetryExhaustedError is the hard, terminal failure retryTransport returns once
// it gives up. It deliberately omits Unwrap into *TransportError so
// errors.AsType can never re-classify it as retryable; Cause stays a plain
// field for logging only.
type RetryExhaustedError struct {
	Failures int
	Cause    error
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("gave up after %d sequential failures: %v", e.Failures, e.Cause)
}

// Source is everything the gate reads from Grafana. httpSource is the one
// production implementation; the tests use a scripted fake
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
// verified against: >= 13.0.0, < 14.0.0.
var (
	supportedGrafanaMin = grafanaVersion{13, 0, 0}
	supportedGrafanaMax = grafanaVersion{14, 0, 0} // exclusive
)

// CheckGrafanaVersion enforces the supported range. An unparseable or
// out-of-range version is a hard error naming both what was found and what is
// supported: the response schemas this gate parses are only verified against
// that range, and trusting an unverified one is how a deprecation turns into a
// silent misread.
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

// httpSource is the production Source: stdlib net/http only, bearer auth from
// a token supplied at construction (the caller reads it from the environment;
// this type never touches env itself), and manual strict decoding via
// ParseState/ParseDefinitions. The retry limit and backoff parameters are
// struct fields with production defaults set here, not package constants, so a
// test can shrink them without a hook.
type httpSource struct {
	baseURL string
	token   string
	client  *http.Client
	clock   Clock

	maxSequentialFailures int
	backoffBase           time.Duration
	backoffCap            time.Duration
}

// NewHTTPSource builds the production Source. token is never logged and never
// enters an error string — it is used only to set the Authorization header.
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
			// change, and the strict parser already turns a real shape change
			// into a loud per-field error the moment it's visible.
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
// the response's own clock.
type requestResult struct {
	Body       []byte
	ServerDate time.Time     // the response's Date header
	Skew       time.Duration // serverDate - (t_send+t_headers)/2, signed
	SkewBound  time.Duration // (t_headers-t_send)/2 — RTT/2 to the response headers
	// Latency spans t_send through the full body read: the budget check needs
	// the whole poll's wall time (header-only latency would be fail-open). The
	// caller's JSON parse runs outside doRequest; extend here if that ever must
	// be folded in.
	Latency time.Duration
}

// doRequest performs one HTTP GET and classifies the outcome: network failure,
// non-2xx, or body-read failure is retryable (*TransportError); a missing or
// unparseable Date header or a skew beyond SkewHardLimit is a hard error —
// retrying can never fix either, so neither enters the backoff loop.
//
// The Date/skew check runs on every endpoint (even /api/health): a skew only
// noticed once RuleState starts polling has already masked earlier reads, so it
// fails closed on the first response.
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
		return requestResult{}, fmt.Errorf("%s: response has no Date header", path)
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
		return requestResult{}, fmt.Errorf("%s: clock skew %s exceeds hard limit %s", path, absSkew, SkewHardLimit)
	}

	return requestResult{Body: b, ServerDate: serverDate, Skew: signedSkew, SkewBound: bound, Latency: latency}, nil
}

// retryTransport runs fn, retrying with backoff only on *TransportError — any
// other error returns immediately. failures counts consecutive *TransportError
// results; exceeding maxFailures gives up with a wrapped hard error. Waits go
// through clock.After so a fake Clock never sleeps real time.
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
// ±20% jitter.
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
