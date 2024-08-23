package wasp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSmokeSchedules(t *testing.T) {
	type test struct {
		name   string
		input  []*Segment
		output []*Segment
	}

	tests := []test{
		{
			name:  "increasing line",
			input: Steps(10, 10, 3, 30*time.Second),
			output: []*Segment{
				{
					From:     10,
					Duration: 10 * time.Second,
				},
				{
					From:     20,
					Duration: 10 * time.Second,
				},
				{
					From:     30,
					Duration: 10 * time.Second,
				},
			},
		},
		{
			name:  "decreasing line",
			input: Steps(100, -10, 3, 30*time.Second),
			output: []*Segment{
				{
					From:     100,
					Duration: 10 * time.Second,
				},
				{
					From:     90,
					Duration: 10 * time.Second,
				},
				{
					From:     80,
					Duration: 10 * time.Second,
				},
			},
		},
		{
			name:  "plain",
			input: Plain(1, 1*time.Second),
			output: []*Segment{
				{
					From:     1,
					Duration: 1 * time.Second,
				},
			},
		},
		{
			name: "combine",
			input: Combine(
				Plain(200, 1*time.Second),
				Plain(300, 1*time.Second),
			),
			output: []*Segment{
				{
					From:     200,
					Duration: 1 * time.Second,
				},
				{
					From:     300,
					Duration: 1 * time.Second,
				},
			},
		},
		{
			name: "combine and repeat",
			input: CombineAndRepeat(
				2,
				Plain(1, 1*time.Second),
				Plain(100, 1*time.Second),
			),
			output: []*Segment{
				{
					From:     1,
					Duration: 1 * time.Second,
				},
				{
					From:     100,
					Duration: 1 * time.Second,
				},
				{
					From:     1,
					Duration: 1 * time.Second,
				},
				{
					From:     100,
					Duration: 1 * time.Second,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.input, tc.output)
		})
	}
}
