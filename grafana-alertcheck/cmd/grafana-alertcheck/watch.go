package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

const watchUsage = "usage: grafana-alertcheck watch --out <file> [--pidfile F] [--daemon-log F] " +
	"--alerts <file|-> [--folder F] [--poll-interval D] [--concurrency N] [--until RFC3339]"

// runWatch is the record step's entire CLI surface, split in two by one flag
// set — gate.DaemonChildFlag ("--daemon-child") and gate.ReadyFDFlag
// ("--ready-fd") select which side of P6's parent/child split this
// invocation is:
//
//   - without them: the record command an operator types. It parses --out,
//     --alerts and the rest, builds a gate.WatchConfig and calls gate.Watch,
//     which resolves, records the first observation of every rule, and
//     detaches the recorder before returning (§4.3).
//   - with them: the detached recorder itself. gate.Watch's own childArgs
//     (watch_unix.go) is the only thing that ever sets them — an operator
//     never types "--daemon-child" and it does not appear in watchUsage —
//     and this dispatches straight to gate.RunDaemonChild. P6's integration
//     test already covers the spawn; this is the one new test P10 owns: that
//     seeing the flag reaches RunDaemonChild.
//
// Both flags live in the SAME flag set as the operator-facing ones rather
// than a second, hidden set: the child is started with childArgs' exact
// argv, e.g. "watch --daemon-child --out log.jsonl --ready-fd 3
// [--until ...] [--concurrency ...]", and a second parser would have to stay
// byte-for-byte in sync with that slice to accept it.
func runWatch(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprintln(stderr, watchUsage) }

	common := registerCommon(fs)
	out := fs.String("out", "", "JSONL log path to record to")
	pidfile := fs.String("pidfile", "", "pidfile path (default <out>.pid)")
	daemonLog := fs.String("daemon-log", "", "stdout/stderr sink for the detached recorder (default <out>.daemon.log)")
	until := fs.String("until", "", "optional hard stop, RFC3339 (default: run until check stops it)")
	pollInterval := fs.String("poll-interval", "", "override every rule's poll cadence (default: half its own evaluation interval)")

	// Hidden: never in watchUsage, never typed by an operator (see doc comment).
	daemonChild := fs.Bool(gate.DaemonChildFlag[2:], false, "")
	readyFD := fs.Int(gate.ReadyFDFlag[2:], 0, "")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *daemonChild {
		return runDaemonChild(*out, *until, *common.concurrency, *readyFD, stderr)
	}

	url, token, err := grafanaEnv()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	alerts, err := readAlerts(stdin, *common.alerts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	cfg := gate.WatchConfig{
		URL:         url,
		Token:       token,
		Alerts:      alerts,
		Folder:      *common.folder,
		Out:         *out,
		PidFile:     *pidfile,
		DaemonLog:   *daemonLog,
		Concurrency: *common.concurrency,
		Clock:       gate.SystemClock{},
		Notes:       stderr,
	}
	if *until != "" {
		t, err := time.Parse(time.RFC3339, *until)
		if err != nil {
			fmt.Fprintf(stderr, "--until: %v\n", err)
			return 2
		}
		cfg.Until = t
	}
	if *pollInterval != "" {
		d, err := time.ParseDuration(*pollInterval)
		if err != nil {
			fmt.Fprintf(stderr, "--poll-interval: %v\n", err)
			return 2
		}
		cfg.PollEvery = d
	}

	if err := gate.Watch(context.Background(), cfg); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}

// runDaemonChild is the detached recorder's whole entry point (P6's
// obligation on this phase). Its stdout and stderr are already the daemon
// log file — spawnChild (watch_unix.go) redirects both before Start — so
// writing to stderr here lands exactly where waitForChildReady's failure
// path quotes from.
func runDaemonChild(out, until string, concurrency, readyFD int, stderr io.Writer) int {
	url, token, err := grafanaEnv()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	cfg := gate.DaemonChildConfig{
		URL:         url,
		Token:       token,
		Out:         out,
		Concurrency: concurrency,
		Clock:       gate.SystemClock{},
		ReadyFD:     readyFD,
	}
	if until != "" {
		t, err := time.Parse(time.RFC3339, until)
		if err != nil {
			fmt.Fprintf(stderr, "--until: %v\n", err)
			return 2
		}
		cfg.Until = t
	}
	if err := gate.RunDaemonChild(context.Background(), cfg); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}

// The daemon-child and ready-fd flags registered above must keep matching
// watch_unix.go's childArgs, which names exactly --daemon-child, --out,
// --ready-fd, --until and --concurrency and nothing else: that function
// builds this process's own argv when it re-execs itself as the detached
// recorder, so a flag added to one side without the other means the child
// fails on its very first flag.Parse.
