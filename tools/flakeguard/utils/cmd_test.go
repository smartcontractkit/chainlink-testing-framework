package utils

import (
	"strings"
	"testing"
)

func TestExecuteCmd(t *testing.T) {
	tests := []struct {
		name          string
		cmd           string
		args          []string
		expectedOut   string
		expectedErr   string
		expectingFail bool
	}{
		{
			name:        "successful command",
			cmd:         "echo",
			args:        []string{"Hello, world!"},
			expectedOut: "Hello, world!\n",
			expectedErr: "",
		},
		{
			name:        "successful command with no args",
			cmd:         "echo",
			args:        []string{},
			expectedOut: "\n",
			expectedErr: "",
		},
		{
			name:          "nonexistent command",
			cmd:           "nonexistentcommand",
			args:          []string{},
			expectedOut:   "",
			expectingFail: true,
		},
		{
			name:          "command with stderr",
			cmd:           "ls",
			args:          []string{"nonexistentfile"},
			expectedOut:   "",
			expectedErr:   "No such file or directory",
			expectingFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteCmd(tt.cmd, tt.args...)
			if (err != nil) != tt.expectingFail {
				t.Fatalf("expected failure: %v, got error: %v", tt.expectingFail, err)
			}

			if output.Stdout.String() != tt.expectedOut {
				t.Errorf("unexpected stdout: expected %q, got %q", tt.expectedOut, output.Stdout.String())
			}

			if !strings.Contains(output.Stderr.String(), tt.expectedErr) {
				t.Errorf("unexpected stderr: expected to contain %q, got %q", tt.expectedErr, output.Stderr.String())
			}
		})
	}
}
