package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

const stopUsage = "usage: grafana-alertcheck stop --out <file> [--pidfile F]"

// runStop reaps a detached recorder without reading or classifying its log —
// what an `if: always()` step calls when the work failed. It is idempotent, so
// it is a no-op after check has already stopped the recorder.
func runStop(args []string, _ io.Reader, _, stderr io.Writer) int {
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprintln(stderr, stopUsage) }

	out := fs.String("out", "", "JSONL log path whose recorder to stop")
	pidfile := fs.String("pidfile", "", "pidfile path (default <out>.pid)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(stderr, "stop: unexpected arguments %v\n", fs.Args())
		return 2
	}
	if *out == "" {
		fmt.Fprintln(stderr, "stop: --out is required")
		return 2
	}

	// SIGINT/SIGTERM cancel the wait cleanly; an interrupted cleanup is a
	// could-not-complete, never a silent success.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	held, err := gate.StopRecorder(ctx, gate.StopConfig{
		Log:     *out,
		PidFile: *pidfile,
		Clock:   gate.SystemClock{},
		Notes:   newNoteStyler(stderr),
		Cleanup: true,
	})
	if held != nil {
		_ = held.Close()
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return 0
}
