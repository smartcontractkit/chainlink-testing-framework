package framework

import (
	"io"
	"strings"
	"testing"
)

func TestCheckNodeLogErrorsFromStreams(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "Clean",
			content:     `{"level":"info","msg":"all good"}`,
			expectError: false,
		},
		{
			name:        "Contains CRIT",
			content:     `{"level":"error","msg":"CRIT happened"}`,
			expectError: true,
		},
		{
			name:        "Contains PANIC",
			content:     `PANIC: something bad`,
			expectError: true,
		},
		{
			name:        "Contains FATAL",
			content:     `{"msg":"FATAL condition"}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streams := map[string]io.ReadCloser{
				tt.name: io.NopCloser(strings.NewReader(tt.content)),
			}

			err := checkNodeLogErrorsFromStreams(streams)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}
