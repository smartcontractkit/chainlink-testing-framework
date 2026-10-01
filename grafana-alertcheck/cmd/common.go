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

// commonFlags is registerCommon's result: the flags watch and check share.
// Connection details are never flags, and states / poll-interval are
// deliberately NOT here — states is check-only because recording is
// unfiltered, and poll-interval is watch-only because check reads the cadence
// from the log header. Putting either here would give both commands an opinion
// about a value only one of them may set.
type commonFlags struct {
	folder        *string
	concurrency   *int
	alerts        *string
	excludeAlerts *string
	includeLabels *string
	excludeLabels *string
}

func registerCommon(fs *flag.FlagSet) *commonFlags {
	return &commonFlags{
		folder:      fs.String("folder", "", "default folder to scope an unqualified alert name to"),
		concurrency: fs.Int("concurrency", 1, "maximum concurrent requests to Grafana"),
		alerts:      fs.String("alerts", "", "path to a file of alert names, one per line, or - for stdin"),
		excludeAlerts: fs.String("exclude-alerts", "",
			"path to a file of alert names to subtract from the selected set, one per line, or - for stdin"),
		includeLabels: fs.String("include-labels", "",
			"comma-separated key=value pairs selecting rules by label, e.g. team=bcm,env=stage (cannot be combined with --alerts)"),
		excludeLabels: fs.String("exclude-labels", "",
			"comma-separated key=value pairs; rules carrying any of them are dropped (requires --include-labels)"),
	}
}

// parseLabelPairs parses a comma-separated list of exact-match key=value label
// pairs. An empty string means the flag was not given.
func parseLabelPairs(flagName, s string) ([]gate.LabelMatcher, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	seen := make(map[string]bool)
	var out []gate.LabelMatcher
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("%s: empty label pair in %q", flagName, s)
		}
		// An empty value is legal: it selects rules that carry the label with
		// an empty value (a missing label never matches).
		key, value, ok := strings.Cut(part, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" {
			return nil, fmt.Errorf("%s: %q is not a key=value pair with a non-empty key", flagName, part)
		}
		if seen[key] {
			return nil, fmt.Errorf("%s: duplicate label %q", flagName, key)
		}
		seen[key] = true
		out = append(out, gate.LabelMatcher{Key: key, Value: value})
	}
	return out, nil
}

// readAlerts reads alert names, one per line, from a file or stdin ("-").
// flagName names the caller's flag so errors point at the right input. An
// empty path is not an error: callers decide whether an empty list is allowed.
func readAlerts(stdin io.Reader, flagName, path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	var r io.Reader
	if path == "-" {
		r = stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("read %s %s: %w", flagName, path, err)
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
		return nil, fmt.Errorf("read %s %s: %w", flagName, path, err)
	}
	return lines, nil
}

// parseStates parses check's --states flag: a comma-separated list of the
// "bad" state vocabulary Config.States matches against (classify.go's
// badStateSet). An empty string is not resolved here — it means "use the
// library default of {firing}" — so this returns nil, nil for "" rather than
// an error.
//
// normal is deliberately NOT accepted. The vocabulary is fixed to
// firing | pending | nodata | error precisely because "normal" is the good
// state, never a bad one to classify against: --states normal would turn every
// healthy instance into a violation and fail every healthy fleet.
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

// parsePreexisting parses check's --preexisting flag.
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
