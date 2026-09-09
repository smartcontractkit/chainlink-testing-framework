package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestMain doubles this test binary as the detached recorder. Watch spawns
// os.Executable(), which under `go test` is this binary, so the one integration
// test below exercises the real thing — a real fork/exec, a real setsid, a real
// inherited environment, a real SIGTERM — with this function standing in for
// the CLI's `watch --daemon-child` dispatch.
func TestMain(m *testing.M) {
	if path := os.Getenv(lockHolderEnv); path != "" {
		os.Exit(runTestLockHolder(path))
	}
	if slices.Contains(os.Args, DaemonChildFlag) {
		os.Exit(runTestDaemonChild(os.Args[1:]))
	}
	os.Exit(m.Run())
}

// lockHolderEnv turns this test binary into a stand-in recorder that holds the
// log's flock and refuses to die: a process check's stop protocol must wait
// for and, on a timeout, refuse to read around.
//
// It has to be a real second process. flock is what stopRecorder probes, and
// there is no flock(1) on darwin, so a shell one-liner cannot take the lock —
// while a lock taken in the test process itself would be granted to the probe
// on some platforms and prove nothing.
const lockHolderEnv = "GRAFANA_ALERTCHECK_TEST_LOCK_HOLDER"

// runTestLockHolder takes the log's exclusive lock, reports that it has it on
// stdout, ignores SIGTERM, and waits to be killed. The report is what lets the
// test start only once the lock is genuinely held, rather than racing it.
func runTestLockHolder(path string) int {
	signal.Ignore(syscall.SIGTERM)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := lockExclusive(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println("locked")

	// Long enough to outlive any test that starts it; the test kills it, and
	// SIGKILL is not ignorable.
	time.Sleep(5 * time.Minute)
	return 0
}

// runTestDaemonChild parses the child argv childArgs() writes, and reads the
// connection details from the environment — never from argv. The CLI's `watch`
// FlagSet does the same four flags.
func runTestDaemonChild(args []string) int {
	cfg := DaemonChildConfig{
		URL:   os.Getenv("GRAFANA_URL"),
		Token: os.Getenv("GRAFANA_TOKEN"),
	}
	for i := 0; i < len(args); i++ {
		value := func() string {
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "flag %s wants a value\n", args[i])
				os.Exit(2)
			}
			i++
			return args[i]
		}
		switch args[i] {
		case "--out":
			cfg.Out = value()
		case "--until":
			until, err := time.Parse(time.RFC3339, value())
			if err != nil {
				fmt.Fprintf(os.Stderr, "--until: %v\n", err)
				return 2
			}
			cfg.Until = until
		case "--concurrency":
			n, err := strconv.Atoi(value())
			if err != nil {
				fmt.Fprintf(os.Stderr, "--concurrency: %v\n", err)
				return 2
			}
			cfg.Concurrency = n
		case ReadyFDFlag:
			fd, err := strconv.Atoi(value())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", ReadyFDFlag, err)
				return 2
			}
			cfg.ReadyFD = fd
		}
	}
	if err := RunDaemonChild(context.Background(), cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	return 0
}

// testBearerToken is what every request to grafanaTestServer must carry. The
// child never receives it in argv, so a request that arrives
// authenticated is proof that the token reached the detached process through
// the inherited environment — and a 401 is what a test sees if that ever
// breaks.
const testBearerToken = "test-token"

// grafanaTestServer serves the three endpoints the record step reads, from the
// real captured fixtures: /api/health, the ruler definitions, and the
// rule_name-filtered state response. The state body is the one-instance
// fixture with its uid and name patched to the ruler fixture's live rule, so
// the reducer's select-by-UID finds it.
func grafanaTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	ruler := readFixture(t, "ruler_rules.json")
	state := patchedStateBody(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+testBearerToken {
			// Not t.Errorf: this must reach the client as a real 401, so the
			// parent fails its version gate and the child fails its polls.
			http.Error(w, "unauthorized: Authorization = "+got, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/health":
			fmt.Fprint(w, healthBody("13.1.0"))
		case strings.HasPrefix(r.URL.Path, "/api/ruler/"):
			_, _ = w.Write(ruler)
		case strings.HasPrefix(r.URL.Path, "/api/prometheus/"):
			if r.URL.Query().Get("rule_name") == "" {
				// The gate must never read the state endpoint unfiltered.
				http.Error(w, "unfiltered state read", http.StatusBadRequest)
				return
			}
			_, _ = w.Write(state)
		default:
			http.Error(w, "unexpected path "+r.URL.Path, http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func patchedStateBody(t *testing.T) []byte {
	t.Helper()
	var body map[string]any
	require.NoError(t, json.Unmarshal(readFixture(t, "state_one_instance.json"), &body))
	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "state fixture: no data object")
	groups, ok := data["groups"].([]any)
	require.True(t, ok, "state fixture: no groups")
	require.NotEmpty(t, groups, "state fixture: no groups")
	group, ok := groups[0].(map[string]any)
	require.True(t, ok, "state fixture: group 0 is not an object")
	rules, ok := group["rules"].([]any)
	require.True(t, ok, "state fixture: group 0 has no rules")
	require.NotEmpty(t, rules, "state fixture: group 0 has no rules")
	rule, ok := rules[0].(map[string]any)
	require.True(t, ok, "state fixture: rule 0 is not an object")
	rule["uid"] = watchActiveUID
	rule["name"] = watchActiveTitle
	rule["lastEvaluation"] = time.Now().UTC().Format(time.RFC3339Nano)

	b, err := json.Marshal(body)
	require.NoError(t, err)
	return b
}

// waitFor polls cond until it holds. This is the one tier of the project where
// a test waits on real time: it drives real processes over real HTTP, so there
// is no clock to fake.
func waitFor(t *testing.T, what string, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	require.Fail(t, fmt.Sprintf("timed out after %s waiting for %s", timeout, what))
}

// The one watch integration test: everything from the version gate to the
// sentinel, through a real detached process.
//
// It asserts the four things only a real spawn can show — the pidfile points
// at a live process, that process is in its own session (setsid, not a bare
// `&`), it keeps appending after Watch returned, and SIGTERM makes it finish
// the log in the stop order — and it uses a 200ms --poll-interval to do it in
// about a second, which also exercises the unclamped-override path.
func TestWatchSpawnsADetachedRecorder(t *testing.T) {
	srv := grafanaTestServer(t)
	t.Setenv("GRAFANA_URL", srv.URL)
	t.Setenv("GRAFANA_TOKEN", testBearerToken)

	out := filepath.Join(t.TempDir(), "log.jsonl")
	var notes strings.Builder
	cfg := WatchConfig{
		URL:         srv.URL,
		Token:       testBearerToken,
		Alerts:      []string{"uid:" + watchActiveUID},
		Out:         out,
		PollEvery:   200 * time.Millisecond,
		Concurrency: 2,
		Notes:       &notes,
	}

	require.NoError(t, Watch(context.Background(), cfg))
	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("notes:\n%s", notes.String())
			t.Logf("daemon log:\n%s", daemonLogTail(out+".daemon.log", 0))
		}
	})

	pid, err := ReadPidFile(out + ".pid")
	require.NoError(t, err)
	require.NoError(t, syscall.Kill(pid, 0), "recorder pid %d is not running right after Watch returned", pid)
	// Setsid, not a bare `&`: a session leader's process group id is its own
	// pid. Without this the child would still share the parent's process group
	// and die with the step that started it.
	pgid, err := syscall.Getpgid(pid)
	require.NoError(t, err)
	require.Equal(t, pid, pgid, "it did not get its own session")

	// The parent already wrote the first heartbeat before it returned;
	// these later ones prove the detached child is the one appending now.
	waitFor(t, "the detached recorder to append its own polls", 10*time.Second, func() bool {
		_, polls, _, err := ReadLog(out)
		return err == nil && len(polls) >= 3
	})

	// Stop it exactly the way check does.
	require.NoError(t, syscall.Kill(pid, syscall.SIGTERM))
	waitFor(t, "the stopped sentinel", 10*time.Second, func() bool {
		_, _, sentinel, err := ReadLog(out)
		return err == nil && sentinel != nil
	})

	header, polls, sentinel, err := ReadLog(out)
	require.NoError(t, err)
	require.Equal(t, srv.URL, header.URL)
	require.Equal(t, "13.1.0", header.GrafanaVersion)
	require.Len(t, header.Rules, 1)
	require.Equal(t, float64(0.2), header.Rules[0].PollEverySeconds)
	for i, p := range polls {
		require.Equalf(t, watchActiveUID, p.RuleUID, "poll %d", i)
		require.Truef(t, p.Found, "poll %d", i)
		require.Falsef(t, p.GrafanaNow.IsZero(), "poll %d has no grafana_now; every poll needs the Date header of its own response", i)
	}
	require.False(t, sentinel.Before(header.StartedAt), "sentinel precedes the record start")

	waitFor(t, "the recorder to exit", 10*time.Second, func() bool {
		return syscall.Kill(pid, 0) != nil
	})
}

// TestDaemonChildRejectsAnAlreadyFinishedLog covers the RunDaemonChild guard
// against a reused --out path: a log that already carries a stopped sentinel is
// a finished recording, and a child starting against it would either append to
// a window already declared over, or take a flock over evidence that is about
// to be classified — so it must refuse before polling once. This is the
// fail-closed counterpart of §4.5 on the recorder's own startup path.
func TestDaemonChildRejectsAnAlreadyFinishedLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	clock := newFakeClock(testNow)

	w, err := NewWriter(path, clock)
	require.NoError(t, err)
	require.NoError(t, w.WriteHeader(testHeader()))
	require.NoError(t, w.Stop())

	err = RunDaemonChild(context.Background(), DaemonChildConfig{
		URL:   testHeader().URL,
		Out:   path,
		Clock: clock,
	})
	require.Error(t, err, "no error against a log that already carries a stopped sentinel")
	require.Contains(t, err.Error(), "sentinel")
}

// TestWatchFailsWhenTheChildCannotStartRecording is the other half of the
// readiness contract. The child dies on its identity check, so it never reports
// ready — and Watch must say so instead of returning success over a window
// nothing is recording, and must leave no pidfile naming a dead process for the
// next step to signal.
//
// The child is made to fail through the environment, which is the only channel
// it takes its connection details from: the header says one URL and the
// inherited GRAFANA_URL says another.
func TestWatchFailsWhenTheChildCannotStartRecording(t *testing.T) {
	srv := grafanaTestServer(t)
	t.Setenv("GRAFANA_URL", srv.URL+"/somewhere-else")
	t.Setenv("GRAFANA_TOKEN", testBearerToken)

	out := filepath.Join(t.TempDir(), "log.jsonl")
	var notes strings.Builder
	err := Watch(context.Background(), WatchConfig{
		URL:         srv.URL, // what the parent uses, and what the header records
		Token:       testBearerToken,
		Alerts:      []string{"uid:" + watchActiveUID},
		Out:         out,
		PollEvery:   200 * time.Millisecond,
		Concurrency: 2,
		Notes:       &notes,
	})
	require.Error(t, err, "the child could never have started recording")
	require.Contains(t, err.Error(), "records url")
	_, statErr := os.Stat(out + ".pid")
	require.True(t, os.IsNotExist(statErr),
		"a pidfile survived a failed detach; pids are reused, so the next step would signal a stranger")
}
