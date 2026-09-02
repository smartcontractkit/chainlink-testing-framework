package gate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParsePromDuration(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"", 0},
		{"0", 0},
		{"0s", 0},
		{"500ms", 500 * time.Millisecond},
		{"1s", time.Second},
		{"1m", time.Minute},
		{"2m", 2 * time.Minute},
		{"3m", 3 * time.Minute},
		{"5m", 5 * time.Minute},
		{"10m", 10 * time.Minute},
		{"15m", 15 * time.Minute},
		{"20m", 20 * time.Minute},
		{"30m", 30 * time.Minute},
		{"1h", time.Hour},
		{"6h", 6 * time.Hour},
		{"12h", 12 * time.Hour},
		{"1d", 24 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1y", 365 * 24 * time.Hour},
		{"1h30m", time.Hour + 30*time.Minute},
		{"1s500ms", time.Second + 500*time.Millisecond},
		{"2d12h", 2*24*time.Hour + 12*time.Hour},
	}
	for _, c := range cases {
		got, err := ParsePromDuration(c.in)
		require.NoErrorf(t, err, "ParsePromDuration(%q)", c.in)
		require.Equalf(t, c.want, got, "ParsePromDuration(%q)", c.in)
	}
}

func TestParsePromDuration_Errors(t *testing.T) {
	cases := []string{
		"5",                   // bare number, no unit
		"-5m",                 // negative
		"5x",                  // unknown unit
		"30m1h",               // ascending order (must be descending)
		"1h1h",                // duplicate unit
		"m",                   // unit with no number
		"1.5h",                // fractional number not supported by this grammar
		"1 h",                 // whitespace
		"300y",                // overflows time.Duration (int64 nanoseconds) — must error, not wrap negative
		"1w2d3h4m5s6ms7us8ns", // too many units
		"carrot",              // completely invalid
	}
	for _, in := range cases {
		_, err := ParsePromDuration(in)
		require.Errorf(t, err, "ParsePromDuration(%q): expected an error, got none", in)
	}
}
