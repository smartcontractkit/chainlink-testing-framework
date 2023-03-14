package loadgen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSchedules(t *testing.T) {
	type test struct {
		name   string
		input  []*Segment
		output []*Segment
	}

	tests := []test{
		{
			name:  "increasing line",
			input: Line(1, 100, 1*time.Second),
			output: []*Segment{
				{
					From:         1,
					Increase:     10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
			},
		},
		{
			name:  "decreasing line",
			input: Line(10, 0, 1*time.Second),
			output: []*Segment{
				{
					From:         10,
					Increase:     -1,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
			},
		},
		{
			name: "combine lines",
			input: Combine(
				Line(1, 100, 1*time.Second),
				Plain(200, 1*time.Second),
				Line(100, 1, 1*time.Second),
			),
			output: []*Segment{
				{
					From:         1,
					Increase:     10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         200,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         100,
					Increase:     -10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
			},
		},
		{
			name: "combine disjointed lines",
			input: Combine(
				Line(1, 100, 1*time.Second),
				Line(1, 300, 1*time.Second),
			),
			output: []*Segment{
				{
					From:         1,
					Increase:     10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         1,
					Increase:     30,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
			},
		},
		{
			name: "combine and repeat",
			input: CombineAndRepeat(
				2,
				Line(1, 100, 1*time.Second),
				Line(100, 1, 1*time.Second),
			),
			output: []*Segment{
				{
					From:         1,
					Increase:     10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         100,
					Increase:     -10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         1,
					Increase:     10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
				},
				{
					From:         100,
					Increase:     -10,
					Steps:        DefaultStepChangePrecision,
					StepDuration: 100 * time.Millisecond,
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
