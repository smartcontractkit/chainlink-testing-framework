package gate

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// promDurationUnit is one accepted unit in a Prometheus-style duration string.
// rank increases with unit size; ParsePromDuration requires strictly decreasing
// rank across concatenated components (e.g. "1h30m", never "30m1h").
type promDurationUnit struct {
	suffix string
	mult   time.Duration
	rank   int
}

// Longest suffix must be tried first ("ms" before "m") — see the matching loop below.
var promDurationUnits = []promDurationUnit{
	{"ms", time.Millisecond, 0},
	{"s", time.Second, 1},
	{"m", time.Minute, 2},
	{"h", time.Hour, 3},
	{"d", 24 * time.Hour, 4},
	{"w", 7 * 24 * time.Hour, 5},
	{"y", 365 * 24 * time.Hour, 6},
}

// ParsePromDuration parses a Grafana/Prometheus-style duration ("1h30m", "1d", "1w").
// Unlike time.ParseDuration, it accepts "d" and "w". "" and "0" are 0.
func ParsePromDuration(s string) (time.Duration, error) {
	if s == "" || s == "0" {
		return 0, nil
	}
	if strings.HasPrefix(s, "-") {
		return 0, fmt.Errorf("invalid duration %q: negative durations are not supported", s)
	}

	var total time.Duration
	prevRank := len(promDurationUnits) // sentinel higher than any real rank
	rest := s
	for rest != "" {
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("invalid duration %q: expected a number", s)
		}
		numPart := rest[:i]
		rest = rest[i:]

		matched := -1
		matchLen := 0
		for idx, u := range promDurationUnits {
			if strings.HasPrefix(rest, u.suffix) && len(u.suffix) > matchLen {
				matched = idx
				matchLen = len(u.suffix)
			}
		}
		if matched == -1 {
			return 0, fmt.Errorf("invalid duration %q: unrecognized unit", s)
		}
		u := promDurationUnits[matched]
		if u.rank >= prevRank {
			return 0, fmt.Errorf("invalid duration %q: units must appear in descending order", s)
		}
		prevRank = u.rank

		n, err := strconv.ParseInt(numPart, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", s, err)
		}
		// n and u.mult are both non-negative here (the leading "-" check
		// above rejects negative input), so overflow of either the
		// multiplication or the running sum can only wrap upward past
		// math.MaxInt64 — check both explicitly rather than let a duration
		// like "300y" silently become negative garbage that would later
		// feed transitionGrace.
		if n != 0 && n > math.MaxInt64/int64(u.mult) {
			return 0, fmt.Errorf("invalid duration %q: overflows time.Duration", s)
		}
		delta := time.Duration(n) * u.mult
		if total > math.MaxInt64-delta {
			return 0, fmt.Errorf("invalid duration %q: overflows time.Duration", s)
		}
		total += delta
		rest = rest[matchLen:]
	}
	return total, nil
}
