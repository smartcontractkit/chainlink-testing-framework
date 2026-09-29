package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/internal/gate"
)

func TestParseLabelPairs(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []gate.LabelMatcher
		wantErr string
	}{
		{"empty means not given", "", nil, ""},
		{"whitespace means not given", "   ", nil, ""},
		{"single pair", "team=bcm", []gate.LabelMatcher{{Key: "team", Value: "bcm"}}, ""},
		{"spaces are trimmed", " team = bcm , env = stage ", []gate.LabelMatcher{
			{Key: "team", Value: "bcm"},
			{Key: "env", Value: "stage"},
		}, ""},
		{"value may contain equals", "query=a=b", []gate.LabelMatcher{{Key: "query", Value: "a=b"}}, ""},
		{"missing equals", "team", nil, "not a non-empty key=value pair"},
		{"empty key", "=bcm", nil, "not a non-empty key=value pair"},
		{"empty value", "team=", nil, "not a non-empty key=value pair"},
		{"empty segment", "team=bcm,", nil, "empty label pair"},
		{"duplicate key", "team=a,team=b", nil, "duplicate label"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLabelPairs("--include-labels", tt.in)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
