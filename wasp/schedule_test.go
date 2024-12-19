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
					Type:     SegmentType_Steps,
				},
				{
					From:     20,
					Duration: 10 * time.Second,
					Type:     SegmentType_Steps,
				},
				{
					From:     30,
					Duration: 10 * time.Second,
					Type:     SegmentType_Steps,
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
					Type:     SegmentType_Steps,
				},
				{
					From:     90,
					Duration: 10 * time.Second,
					Type:     SegmentType_Steps,
				},
				{
					From:     80,
					Duration: 10 * time.Second,
					Type:     SegmentType_Steps,
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
					Type:     SegmentType_Plain,
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
					Type:     SegmentType_Plain,
				},
				{
					From:     300,
					Duration: 1 * time.Second,
					Type:     SegmentType_Plain,
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
					Type:     SegmentType_Plain,
				},
				{
					From:     100,
					Duration: 1 * time.Second,
					Type:     SegmentType_Plain,
				},
				{
					From:     1,
					Duration: 1 * time.Second,
					Type:     SegmentType_Plain,
				},
				{
					From:     100,
					Duration: 1 * time.Second,
					Type:     SegmentType_Plain,
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
