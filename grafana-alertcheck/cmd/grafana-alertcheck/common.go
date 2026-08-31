package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

// commonFlags is registerCommon's result: the exactly three flags watch and
// check share (§20). Connection details are never flags (§20.2) and states /
// poll-interval are deliberately NOT here — states is check-only because
// recording is unfiltered (P6), and poll-interval is watch-only because check
// reads the cadence from the log header (P5). Putting either here would
// silently reinstate a knob this plan removed.
type commonFlags struct {
	folder      *string
	concurrency *int
	alerts      *string
}

func registerCommon(fs *flag.FlagSet) *commonFlags {
	return &commonFlags{
		folder:      fs.String("folder", "", "default folder to scope an unqualified alert name to (§17)"),
		concurrency: fs.Int("concurrency", 1, "maximum concurrent requests to Grafana"),
		alerts:      fs.String("alerts", "", "path to a file of alert names, one per line, or - for stdin"),
	}
}

// readAlerts reads §17's alert names, one per line, from a file or from
// stdin when path is "-". An empty path is not an error here — watch and
// check each decide for themselves whether an empty list is allowed
// (log mode never wants one; single-step / record mode always does).
func readAlerts(stdin io.Reader, path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	var r io.Reader
	if path == "-" {
		r = stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("read --alerts %s: %w", path, err)
		}
		defer f.Close()
		r = f
	}
	var lines []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read --alerts %s: %w", path, err)
	}
	return lines, nil
}

// parseStates parses check's --states flag: a comma-separated list of the
// "bad" state vocabulary Config.States matches against (§13, classify.go's
// badStateSet). An empty string is not resolved here — it means "use the
// library default of {firing}" — so this returns nil, nil for "" rather than
// an error.
//
// normal is deliberately NOT accepted: the v2 plan fixes this vocabulary to
// firing | pending | nodata | error (line 378) precisely because "normal" is
// the good state, never a bad one to classify against. Accepting it here
// would let --states normal turn every healthy instance into a violation and
// fail every healthy fleet — the exact fail-open shape H7 exists to prevent.
func parseStates(s string) ([]gate.State, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var out []gate.State
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch gate.State(part) {
		case gate.StateFiring, gate.StatePending, gate.StateNodata, gate.StateError:
			out = append(out, gate.State(part))
		default:
			return nil, fmt.Errorf("--states: unknown state %q (want any of: firing, pending, nodata, error)", part)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("--states: %q named no state", s)
	}
	return out, nil
}

// parsePreexisting parses check's --preexisting flag (§11.7).
func parsePreexisting(s string) (gate.PreexistingPolicy, error) {
	switch gate.PreexistingPolicy(s) {
	case "":
		return gate.PreexistingFailUnlessRecovered, nil
	case gate.PreexistingFailUnlessRecovered, gate.PreexistingFail, gate.PreexistingIgnore:
		return gate.PreexistingPolicy(s), nil
	default:
		return "", fmt.Errorf("--preexisting: unknown policy %q (want one of: fail-unless-recovered, fail, ignore)", s)
	}
}
