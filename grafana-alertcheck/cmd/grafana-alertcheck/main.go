// Command grafana-alertcheck is the CLI entry point for the gate. P3 wires
// only the `list` subcommand — enough to validate auth, the ruler parse, and
// resolution against a real Grafana before any coverage logic exists (§9 rule
// 4, "reach runnable at PR 5"). P10 extends this file with `watch` and
// `check`.
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

const usage = "usage: grafana-alertcheck <list>"

// run is the whole of main's testable surface: parse the subcommand, dispatch,
// return the process exit code. Exit codes below 2 (pass/violations) belong to
// `check` alone (§20.3, P10); every failure reachable from here — a missing
// subcommand, a bad flag, a transport or auth failure — is a could-not-check
// condition and maps to 2, never to 0 or 1 (H7).
//
// Requested help (-h/--help) is not a failure — it is the one exception to
// that rule. Convention (and every stdlib flag.FlagSet default) is exit 0 to
// stdout for help the caller asked for, reserving 2/stderr for help printed
// *because* something else went wrong (no subcommand, an unknown one).
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return 2
	}

	switch args[0] {
	case "list":
		return runList(args[1:], stdout, stderr)
	case "-h", "-help", "--help":
		fmt.Fprintln(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown subcommand %q\n", args[0])
		return 2
	}
}
