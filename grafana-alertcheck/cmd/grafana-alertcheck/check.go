package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os/signal"
	"syscall"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

const checkUsage = "usage: grafana-alertcheck check [--in <file>] [--pidfile F] --from RFC3339 --to RFC3339 " +
	"[--alerts ...] [--folder F] [--states ...] [--preexisting ...] [--min-observed N] [--allow-paused] " +
	"[--nodata-is-unobservable] [--concurrency N] [--output json]"

// runCheck is the classify step's CLI surface: parse flags into a gate.Config,
// run gate.Check, and translate its (Result, error) into output and an exit
// code. All of the correctness lives in the gate package — this file's only job
// is presentation and the exit-code mapping, which exitCode below keeps as one
// pure function so it can be tested without a network.
func runCheck(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprintln(stderr, checkUsage) }

	common := registerCommon(fs)
	in := fs.String("in", "", "path of a log recorded by watch; empty selects single-step mode")
	pidfile := fs.String("pidfile", "", "pidfile of the recorder to stop before reading --in (default <in>.pid)")
	from := fs.String("from", "", "the moment the deploy finished, RFC3339 (required with --in)")
	to := fs.String("to", "", "the end of the window to classify, RFC3339 (required)")
	states := fs.String("states", "", "comma-separated bad states to classify against (default: firing)")
	preexisting := fs.String("preexisting", "", "how to judge an instance already bad at `from` (default: fail-unless-recovered)")
	minObserved := fs.Int("min-observed", 0, "minimum rules that must be observed (default: every resolved rule)")
	allowPaused := fs.Bool("allow-paused", false, "do not count a rule paused before the window against --min-observed")
	nodataIsUnobservable := fs.Bool("nodata-is-unobservable", false, "treat a sustained health=nodata as unobservable rather than a note")
	output := fs.String("output", "", `"json" writes the machine-readable Result to stdout in addition to the table; default is the table alone`)

	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *output != "" && *output != "json" {
		fmt.Fprintf(stderr, "--output: unknown value %q (only \"json\" is supported)\n", *output)
		return 2
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
	stateList, err := parseStates(*states)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	preexistingPolicy, err := parsePreexisting(*preexisting)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	cfg := gate.Config{
		URL:                  url,
		Token:                token,
		Alerts:               alerts,
		Folder:               *common.folder,
		States:               stateList,
		Preexisting:          preexistingPolicy,
		MinObserved:          *minObserved,
		AllowPaused:          *allowPaused,
		NodataIsUnobservable: *nodataIsUnobservable,
		Log:                  *in,
		PidFile:              *pidfile,
		Concurrency:          *common.concurrency,
		Clock:                gate.SystemClock{},
		Notes:                stderr,
	}
	if *to == "" {
		fmt.Fprintln(stderr, "check: --to is required")
		return 2
	}
	t, err := time.Parse(time.RFC3339, *to)
	if err != nil {
		fmt.Fprintf(stderr, "--to: %v\n", err)
		return 2
	}
	cfg.To = t
	if *from != "" {
		f, err := time.Parse(time.RFC3339, *from)
		if err != nil {
			fmt.Fprintf(stderr, "--from: %v\n", err)
			return 2
		}
		cfg.From = f
	}

	// SIGINT/SIGTERM cancel the run cleanly rather than leaving an operator's
	// Ctrl-C to kill the process mid-collection: Check's collection loop and
	// drain wait both already select on ctx.Done() (check.go), so this makes
	// an interrupted run fail the way every other could-not-check path does
	// — exit 2, never a silently truncated pass.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	result, checkErr := gate.Check(ctx, cfg)

	if err := renderTable(stderr, result); err != nil {
		fmt.Fprintln(stderr, err)
	}
	if checkErr != nil {
		fmt.Fprintln(stderr, checkErr)
	}
	if *output == "json" {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			fmt.Fprintf(stderr, "encode --output json: %v\n", err)
			return 2
		}
	}
	return exitCode(result, checkErr)
}

// exitCode is the whole exit-code mapping, kept as one pure function of
// exactly what Check returns so it is testable without a network: err != nil
// is exit 2 UNCONDITIONALLY — never 0 and never 1, even alongside real
// violations, because an inability to check beats a violation and an error is
// never a pass. Violations without an error is exit 1. Neither is exit 0.
func exitCode(res gate.Result, err error) int {
	switch {
	case err != nil:
		return 2
	case len(res.Violations) > 0:
		return 1
	default:
		return 0
	}
}
