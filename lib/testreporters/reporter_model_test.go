package testreporters_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
)

func TestVerifyLogFile(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		failingLogLevel  zapcore.Level
		failureThreshold uint
		allowedMessages  []testreporters.AllowedLogMessage
		expectedError    string
	}{
		{
			name:             "No logs",
			content:          "",
			failingLogLevel:  zapcore.ErrorLevel,
			failureThreshold: 1,
			allowedMessages:  nil,
			expectedError:    "",
		},
		{
			name:             "Log level below threshold",
			content:          `{"level":"error","msg":"info log"}`,
			failingLogLevel:  zapcore.WarnLevel,
			failureThreshold: 2,
			allowedMessages:  nil,
			expectedError:    "",
		},
		{
			name:             "Log level equals threshold with failure",
			content:          `{"level":"error","msg":"error log"}`,
			failingLogLevel:  zapcore.ErrorLevel,
			failureThreshold: 1,
			allowedMessages:  nil,
			expectedError:    "found log at level 'error', failing any log level higher than error: {\"level\":\"error\",\"msg\":\"error log\"}",
		},
		{
			name:             "Log level above threshold with failure",
			content:          `{"level":"error","msg":"critical error"}`,
			failingLogLevel:  zapcore.WarnLevel,
			failureThreshold: 1,
			allowedMessages:  nil,
			expectedError:    "found log at level 'error', failing any log level higher than error: {\"level\":\"error\",\"msg\":\"critical error\"}",
		},
		{
			name: "Log level above threshold with failure",
			content: `
{"level":"error","msg":"critical error 1"}
{"level":"error","msg":"critical error 2"}`,
			failingLogLevel:  zapcore.WarnLevel,
			failureThreshold: 2,
			allowedMessages:  nil,
			expectedError:    "found too many logs at level 'error' or above; failure threshold of 2 reached; last error found: {\"level\":\"error\",\"msg\":\"critical error 2\"}",
		},
		{
			name:             "Allowed message prevents error",
			content:          `{"level":"error","msg":"expected error"}`,
			failingLogLevel:  zapcore.WarnLevel,
			failureThreshold: 1,
			allowedMessages: []testreporters.AllowedLogMessage{
				testreporters.NewAllowedLogMessage("expected", "test", zapcore.ErrorLevel, testreporters.WarnAboutAllowedMsgs_Yes),
			},
			expectedError: "",
		},
		{
			name:             "Threshold set to zero",
			content:          `{"level":"error","msg":"error log"}`,
			failingLogLevel:  zapcore.ErrorLevel,
			failureThreshold: 0, // This should work the same as when it's set to 1
			allowedMessages:  nil,
			expectedError:    "found log at level 'error', failing any log level higher than error: {\"level\":\"error\",\"msg\":\"error log\"}",
		},
		{
			name: "Log level above threshold with failure with multiple allowed messages",
			content: `
{"level":"error","msg":"ignored error"}
{"level":"panic","msg":"ignored critical"}
{"level":"error","msg":"critical error that should fail"}
{"level":"error","msg":"critical error that should fail"}`,
			failingLogLevel:  zapcore.ErrorLevel,
			failureThreshold: 2,
			allowedMessages: []testreporters.AllowedLogMessage{
				testreporters.NewAllowedLogMessage("ignored error", "test", zapcore.ErrorLevel, testreporters.WarnAboutAllowedMsgs_Yes),
				testreporters.NewAllowedLogMessage("ignored critical", "test", zapcore.DPanicLevel, testreporters.WarnAboutAllowedMsgs_Yes),
			},
			expectedError: "found too many logs at level 'error' or above; failure threshold of 2 reached; last error found: {\"level\":\"error\",\"msg\":\"critical error that should fail\"}",
		},
		{
			name: "Log level above threshold with allowed message that should not be logged",
			content: `
{"level":"error","msg":"ignored error"}
{"level":"error","msg":"critical error that should fail"}
{"level":"error","msg":"critical error that should fail"}`,
			failingLogLevel:  zapcore.ErrorLevel,
			failureThreshold: 2,
			allowedMessages: []testreporters.AllowedLogMessage{
				testreporters.NewAllowedLogMessage("ignored error", "test", zapcore.ErrorLevel, testreporters.WarnAboutAllowedMsgs_No),
			},
			expectedError: "found too many logs at level 'error' or above; failure threshold of 2 reached; last error found: {\"level\":\"error\",\"msg\":\"critical error that should fail\"}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "log")
			require.NoError(t, err, "failed to create temporary file")
			defer func() {
				_ = os.Remove(tempFile.Name())
				_ = tempFile.Close()
			}()

			_, err = tempFile.WriteString(tc.content)
			require.NoError(t, err, "Failed to write content to file")
			require.NoError(t, tempFile.Sync(), "Failed to flush writes to storage")
			_, err = tempFile.Seek(0, 0) // Rewind file to the beginning for reading
			require.NoError(t, err, "Failed to move cursor to the beginning of the file")

			err = testreporters.VerifyLogFile(tempFile, tc.failingLogLevel, tc.failureThreshold, tc.allowedMessages...)

			if tc.expectedError == "" {
				require.NoError(t, err, "Expected no error but got one.")
			} else {
				require.EqualError(t, err, tc.expectedError, "Expected error did not match actual error.")
			}
		})
	}
}
