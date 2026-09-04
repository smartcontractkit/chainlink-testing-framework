package gate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func healthBody(version string) string {
	return fmt.Sprintf(`{"database":"ok","version":%q,"commit":"abc123"}`, version)
}

func emptyStateBody() string {
	return `{"status":"success","data":{"groups":[]}}`
}

// rawHTTPServer starts an httptest server whose handler hijacks the
// connection and writes exactly the bytes respond returns, bypassing
// net/http's automatic Date-header insertion — the only way to test a
// response with no Date header at all, or a deliberately garbled one.
func rawHTTPServer(t *testing.T, respond func(r *http.Request) []byte) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Errorf("ResponseWriter does not support Hijacker")
			return
		}
		conn, buf, err := hj.Hijack()
		if err != nil {
			t.Errorf("hijack: %v", err)
			return
		}
		defer conn.Close()
		if _, err := buf.Write(respond(r)); err != nil {
			t.Errorf("write raw response: %v", err)
			return
		}
		_ = buf.Flush()
	}))
	t.Cleanup(srv.Close)
	return srv
}

// rawResponse builds a minimal, fully-controlled HTTP/1.1 response: no
// header net/http would add unasked, in particular no automatic Date.
func rawResponse(status int, statusText string, headers map[string]string, body string) []byte {
	out := fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, statusText)
	for k, v := range headers {
		out += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	out += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	out += "Connection: close\r\n\r\n"
	out += body
	return []byte(out)
}

func TestCheckGrafanaVersion(t *testing.T) {
	cases := []struct {
		version      string
		wantErr      bool
		wantContains []string // required substrings of the error message, per case, when wantErr
	}{
		{"13.1.0", false, nil},
		{"13.0.0", false, nil},
		{"13.99.99", false, nil},
		{"12.9.9", true, []string{`"12.9.9"`, "13.0.0", "14.0.0"}},
		{"14.0.0", true, []string{`"14.0.0"`, "13.0.0", "14.0.0"}},
		{"14.1.0", true, []string{`"14.1.0"`, "13.0.0", "14.0.0"}},
		{"not-a-version", true, []string{`"not-a-version"`}},
		{"", true, nil},
	}
	for _, c := range cases {
		err := CheckGrafanaVersion(c.version)
		if c.wantErr && err == nil {
			t.Errorf("CheckGrafanaVersion(%q): want error, got nil", c.version)
			continue
		}
		if !c.wantErr && err != nil {
			t.Errorf("CheckGrafanaVersion(%q): unexpected error: %v", c.version, err)
			continue
		}
		for _, want := range c.wantContains {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("CheckGrafanaVersion(%q): error %q does not mention %q — it must name both what was found and what is supported", c.version, err.Error(), want)
			}
		}
	}
}

func TestBackoffDelay(t *testing.T) {
	base := time.Second
	maxDelay := 30 * time.Second
	maxWithJitter := maxDelay + maxDelay/5 + time.Millisecond
	for n := 1; n <= 10; n++ {
		d := backoffDelay(base, maxDelay, n)
		if d <= 0 {
			t.Fatalf("backoffDelay(_, _, %d) = %v, want > 0", n, d)
		}
		if d > maxWithJitter {
			t.Fatalf("backoffDelay(_, _, %d) = %v, want <= ~%v", n, d, maxWithJitter)
		}
	}
}

func TestHTTPSource_Version_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/health" {
			t.Errorf("path = %q, want /api/health", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(healthBody("13.1.0")))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	v, err := src.Version(context.Background())
	if err != nil {
		t.Fatalf("Version(): unexpected error: %v", err)
	}
	if v != "13.1.0" {
		t.Fatalf("Version() = %q, want 13.1.0", v)
	}
}

func TestHTTPSource_Version_NeverLogsToken(t *testing.T) {
	const secret = "super-secret-token"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+secret {
			t.Errorf("Authorization = %q, want Bearer %s", got, secret)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, secret, clock)
	_, err := src.Version(context.Background())
	if err == nil {
		t.Fatalf("Version(): want error, got nil")
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error %q leaks the token", err.Error())
	}
}

func TestHTTPSource_RuleState_EmptyIsNotAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(emptyStateBody()))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	obs, err := src.RuleState(context.Background(), "Anything")
	if err != nil {
		t.Fatalf("RuleState(): unexpected error: %v", err)
	}
	if len(obs.Rules) != 0 {
		t.Fatalf("Rules = %+v, want empty (an authoritative 2xx is not a transport error)", obs.Rules)
	}
	if obs.GrafanaNow.IsZero() {
		t.Fatalf("GrafanaNow is zero, want the response's Date header value")
	}
}

func TestHTTPSource_RuleState_EscapesRuleName(t *testing.T) {
	var gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		if r.URL.Path != "/api/prometheus/grafana/api/v1/rules" {
			t.Errorf("path = %q, want /api/prometheus/grafana/api/v1/rules", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(emptyStateBody()))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	title := "[JD] No Job Proposals & More"
	if _, err := src.RuleState(context.Background(), title); err != nil {
		t.Fatalf("RuleState(): unexpected error: %v", err)
	}
	want := "rule_name=" + url.QueryEscape(title)
	if gotQuery != want {
		t.Fatalf("query = %q, want %q", gotQuery, want)
	}
}

func TestHTTPSource_Definitions_HappyPath(t *testing.T) {
	body := readFixture(t, "ruler_rules.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ruler/grafana/api/v1/rules" {
			t.Errorf("path = %q, want /api/ruler/grafana/api/v1/rules", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("query = %q, want none — Definitions reads the ruler API unfiltered", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	defs, err := src.Definitions(context.Background())
	if err != nil {
		t.Fatalf("Definitions(): unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatalf("Definitions(): got 0 definitions from a fixture known to have some")
	}
}

func TestHTTPSource_Skew(t *testing.T) {
	cases := []struct {
		name    string
		drift   time.Duration
		wantErr bool
	}{
		{"0s skew is fine", 0 * time.Second, false},
		{"30s skew is fine", 30 * time.Second, false},
		{"60s skew is fine", 60 * time.Second, false},
		{"61s skew is a hard error", 61 * time.Second, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			anchor := time.Now()
			clock := newFakeClock(anchor)
			var calls atomic.Int32
			srv := rawHTTPServer(t, func(r *http.Request) []byte {
				calls.Add(1)
				date := anchor.Add(c.drift).UTC().Format(http.TimeFormat)
				return rawResponse(200, "OK", map[string]string{
					"Content-Type": "application/json",
					"Date":         date,
				}, healthBody("13.1.0"))
			})
			src := NewHTTPSource(srv.URL, "", clock)
			_, err := src.Version(context.Background())
			if c.wantErr && err == nil {
				t.Fatalf("Version(): want error, got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("Version(): unexpected error: %v", err)
			}
			if c.wantErr && calls.Load() != 1 {
				t.Fatalf("calls = %d, want 1 — a skew hard error must never be retried", calls.Load())
			}
		})
	}
}

func TestHTTPSource_MissingDateHeader(t *testing.T) {
	var calls atomic.Int32
	srv := rawHTTPServer(t, func(r *http.Request) []byte {
		calls.Add(1)
		return rawResponse(200, "OK", map[string]string{
			"Content-Type": "application/json",
		}, healthBody("13.1.0"))
	})
	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	_, err := src.Version(context.Background())
	if err == nil {
		t.Fatalf("Version(): want error, got nil: a missing Date header is a hard error")
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1 — a missing Date header must never be retried", calls.Load())
	}
}

func TestHTTPSource_UnparseableDateHeader(t *testing.T) {
	var calls atomic.Int32
	srv := rawHTTPServer(t, func(r *http.Request) []byte {
		calls.Add(1)
		return rawResponse(200, "OK", map[string]string{
			"Content-Type": "application/json",
			"Date":         "definitely not a date",
		}, healthBody("13.1.0"))
	})
	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	_, err := src.Version(context.Background())
	if err == nil {
		t.Fatalf("Version(): want error, got nil: an unparseable Date header is a hard error")
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1 — an unparseable Date header must never be retried", calls.Load())
	}
}

// TestHTTPSource_ObservationTiming pins the arithmetic behind Observation's
// three derived fields, not just the hard-limit behavior TestHTTPSource_Skew
// already covers: the sign of Skew, SkewBound as exactly RTT/2 to the
// response headers, and Latency as the full send-through-body-read span
// (not just the header round trip). steppingClock advances by a fixed 2s on
// every clock.Now() call, and doRequest calls Now() exactly three times per
// attempt (before send, after headers, after the body read), so the
// arithmetic is exact rather than a real-time approximation.
func TestHTTPSource_ObservationTiming(t *testing.T) {
	cases := []struct {
		name  string
		drift time.Duration
	}{
		{"positive skew", 5 * time.Second},
		{"negative skew", -5 * time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			anchor := time.Now().Truncate(time.Second)
			clock := &steppingClock{now: anchor, step: 2 * time.Second}
			// With step=2s: t_send=anchor, t_headers=anchor+2s, t_bodyRead=anchor+4s.
			// mid = t_send + (t_headers-t_send)/2 = anchor+1s, so serverDate =
			// anchor+1s+drift makes Skew land on exactly `drift`.
			srv := rawHTTPServer(t, func(r *http.Request) []byte {
				date := anchor.Add(time.Second + c.drift).UTC().Format(http.TimeFormat)
				return rawResponse(200, "OK", map[string]string{
					"Content-Type": "application/json",
					"Date":         date,
				}, emptyStateBody())
			})
			src := NewHTTPSource(srv.URL, "", clock)
			obs, err := src.RuleState(context.Background(), "Anything")
			if err != nil {
				t.Fatalf("RuleState(): unexpected error: %v", err)
			}
			if obs.Skew != c.drift {
				t.Errorf("Skew = %v, want %v", obs.Skew, c.drift)
			}
			if obs.SkewBound != time.Second {
				t.Errorf("SkewBound = %v, want 1s (RTT/2 with a 2s round trip to headers)", obs.SkewBound)
			}
			if obs.Latency != 4*time.Second {
				t.Errorf("Latency = %v, want 4s (send through full body read) — not just the 2s header round trip", obs.Latency)
			}
		})
	}
}

// A discriminating regression for "the gate compares staleness against the
// Date header, never the runner's clock". lastEvaluation
// sits 100s behind Grafana's TRUE now (obs.GrafanaNow, from the Date header)
// — under the 120s evalStaleAfter limit — but 130s behind the RUNNER's clock.
// An implementation that leaked the runner's clock into the staleness
// comparison, instead of the Date header, would report a false violation
// here; coverage_test.go's TestProveCoverage_SkewTranslationAtWindowBoundary
// cannot catch that, because it sets LastEvaluation equal to GrafanaNow on
// every poll, making staleness zero regardless of which clock is used.
func TestHTTPSourceStalenessNeverFalsePositiveUnderSkew(t *testing.T) {
	runnerNow := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := newFakeClock(runnerNow)
	const skew = 30 * time.Second // the runner's clock reads 30s ahead of Grafana's
	serverDate := runnerNow.Add(-skew)
	lastEval := serverDate.Add(-100 * time.Second)

	def := Definition{UID: "r1", Title: "Rule One"}
	srv := rawHTTPServer(t, func(r *http.Request) []byte {
		body := fmt.Sprintf(`{"status":"success","data":{"groups":[{"file":"F","name":"G","interval":60,"rules":[`+
			`{"uid":%q,"name":%q,"state":"inactive","health":"ok","isPaused":false,"lastEvaluation":%q}`+
			`]}]}}`, def.UID, def.Title, lastEval.UTC().Format(time.RFC3339))
		return rawResponse(200, "OK", map[string]string{
			"Content-Type": "application/json",
			"Date":         serverDate.UTC().Format(http.TimeFormat),
		}, body)
	})

	src := NewHTTPSource(srv.URL, "", clock)
	obs, err := src.RuleState(context.Background(), def.Title)
	if err != nil {
		t.Fatalf("RuleState(): %v", err)
	}
	if !obs.GrafanaNow.Equal(serverDate) {
		t.Fatalf("GrafanaNow = %s, want the Date header %s, never the runner's clock %s", obs.GrafanaNow, serverDate, runnerNow)
	}
	if len(obs.Rules) != 1 {
		t.Fatalf("Rules = %+v, want exactly one", obs.Rules)
	}

	rt := newRuleTimings(30*time.Second, 60) // evalStaleAfter = 120s
	from := serverDate.Add(-10 * time.Minute)
	to := serverDate
	polls := denseHealthyPolls(def.UID, from, to, 30*time.Second)
	polls[len(polls)-1].LastEvaluation = obs.Rules[0].LastEvaluation // the real, HTTP-sourced value
	sentinel := to

	res := proveCoverage(Header{StartedAt: from.Add(-time.Hour)}, polls, &sentinel, rt, def, from, to, 0)
	if res.Unobservable {
		t.Fatalf("Coverage = %+v, want no violation: 100s behind Grafana's TRUE now is under the 120s limit — "+
			"only a runner-clock leak (skewed +30s here) would push this over", res)
	}
}

func TestHTTPSource_Retry_TransientRecovers(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(healthBody("13.1.0")))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	v, err := src.Version(context.Background())
	if err != nil {
		t.Fatalf("Version(): unexpected error after a transient failure: %v", err)
	}
	if v != "13.1.0" {
		t.Fatalf("Version() = %q, want 13.1.0", v)
	}
	mu.Lock()
	n := calls
	mu.Unlock()
	if n != 3 {
		t.Fatalf("calls = %d, want 3 (2 failures + 1 success)", n)
	}
}

func TestHTTPSource_Retry_ExceedsLimit(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	_, err := src.Version(context.Background())
	if err == nil {
		t.Fatalf("Version(): want error, got nil")
	}
	if n := calls.Load(); n != 6 {
		t.Fatalf("calls = %d, want 6 (maxSequentialFailures=5 tolerates 5, gives up on the 6th)", n)
	}
	assertRetryExhausted(t, err, 6)
}

func TestHTTPSource_RuleState_GarbageBodyRetries(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if n <= 2 {
			// A 2xx with a body that fails ParseState — classified as a
			// transient *TransportError (source.go), not a hard schema
			// break, so it must retry rather than fail immediately.
			_, _ = w.Write([]byte("{not valid json"))
			return
		}
		_, _ = w.Write([]byte(emptyStateBody()))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	obs, err := src.RuleState(context.Background(), "Anything")
	if err != nil {
		t.Fatalf("RuleState(): unexpected error after a transient garbage body: %v", err)
	}
	if len(obs.Rules) != 0 {
		t.Fatalf("Rules = %+v, want empty", obs.Rules)
	}
	mu.Lock()
	n := calls
	mu.Unlock()
	if n != 3 {
		t.Fatalf("calls = %d, want 3 (2 unparseable bodies + 1 valid one)", n)
	}
}

func TestHTTPSource_Definitions_GarbageBodyGivesUp(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{not valid json"))
	}))
	defer srv.Close()

	clock := newFakeClock(time.Now())
	src := NewHTTPSource(srv.URL, "", clock)
	_, err := src.Definitions(context.Background())
	if err == nil {
		t.Fatalf("Definitions(): want error, got nil")
	}
	if n := calls.Load(); n != 6 {
		t.Fatalf("calls = %d, want 6 — a persistently unparseable 2xx body retries like any other transport failure", n)
	}
	assertRetryExhausted(t, err, 6)
}

func TestHTTPSource_NetworkFailureRetries(t *testing.T) {
	// A server that is never listening: every attempt is a network failure,
	// classified as *TransportError, so this exercises the same retry path
	// as a 5xx without needing a real listening socket per failure.
	clock := newFakeClock(time.Now())
	src := NewHTTPSource("http://127.0.0.1:1", "", clock)
	_, err := src.Version(context.Background())
	if err == nil {
		t.Fatalf("Version(): want error, got nil")
	}
	assertRetryExhausted(t, err, 6)
}

// assertRetryExhausted checks the two properties a retry give-up must have:
// it names how many failures it gave up after, and — the regression this
// pins — it is never itself classified as a *TransportError. If it were,
// something one layer up that also retries on *TransportError would treat an
// already-exhausted give-up as retryable again.
func assertRetryExhausted(t *testing.T, err error, wantFailures int) {
	t.Helper()
	var reErr *RetryExhaustedError
	if !errors.As(err, &reErr) {
		t.Fatalf("error %v (%T): want a *RetryExhaustedError", err, err)
	}
	if reErr.Failures != wantFailures {
		t.Errorf("RetryExhaustedError.Failures = %d, want %d", reErr.Failures, wantFailures)
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("gave up after %d", wantFailures)) {
		t.Errorf("error %q does not name the failure count", err.Error())
	}
	if _, ok := errors.AsType[*TransportError](err); ok {
		t.Fatalf("error %v (%T) is classified as *TransportError — an exhausted retry must be a terminal, non-retryable error", err, err)
	}
}

func TestFakeClock(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := newFakeClock(start)
	if !c.Now().Equal(start) {
		t.Fatalf("Now() = %v, want %v", c.Now(), start)
	}
	c.Advance(5 * time.Minute)
	want := start.Add(5 * time.Minute)
	if !c.Now().Equal(want) {
		t.Fatalf("Now() after Advance = %v, want %v", c.Now(), want)
	}

	select {
	case fired := <-c.After(time.Hour):
		if !fired.Equal(want.Add(time.Hour)) {
			t.Fatalf("After fired with %v, want %v", fired, want.Add(time.Hour))
		}
	default:
		t.Fatalf("After(1h) did not fire immediately")
	}
}

func TestFakeSource(t *testing.T) {
	f := newFakeSource()
	f.version = "13.1.0"
	f.defs = []Definition{{UID: "u1", Title: "Rule One"}}

	ctx := context.Background()
	if v, err := f.Version(ctx); err != nil || v != "13.1.0" {
		t.Fatalf("Version() = (%q, %v), want (13.1.0, nil)", v, err)
	}
	if defs, err := f.Definitions(ctx); err != nil || len(defs) != 1 {
		t.Fatalf("Definitions() = (%v, %v), want one definition", defs, err)
	}

	f.script("Rule One", Observation{Rules: []StateRule{{UID: "u1"}}}, nil)
	f.script("Rule One", Observation{}, fmt.Errorf("boom"))
	f.script("Rule One", Observation{Rules: nil}, nil)

	obs, err := f.RuleState(ctx, "Rule One")
	if err != nil || len(obs.Rules) != 1 {
		t.Fatalf("RuleState() call 1 = (%v, %v), want one rule, no error", obs, err)
	}
	if _, err := f.RuleState(ctx, "Rule One"); err == nil {
		t.Fatalf("RuleState() call 2: want the scripted error, got nil")
	}
	obs, err = f.RuleState(ctx, "Rule One")
	if err != nil {
		t.Fatalf("RuleState() call 3: unexpected error: %v", err)
	}
	if obs.Rules != nil {
		t.Fatalf("RuleState() call 3: Rules = %v, want nil (last script entry, then repeats)", obs.Rules)
	}
	obs, err = f.RuleState(ctx, "Rule One")
	if err != nil || obs.Rules != nil {
		t.Fatalf("RuleState() call 4: want the last scripted entry to repeat, got (%v, %v)", obs, err)
	}

	if _, err := f.RuleState(ctx, "Unscripted Rule"); err == nil {
		t.Fatalf("RuleState() for an unscripted title: want an error, got nil")
	}
}
