package gotestevent

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

//nolint:gosec
var inputPassFailMixed = `{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"=== RUN   TestPassTest\n"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"    environment.go:1023: + ./remote-runner.test\n"}
{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"--- PASS: TestPassTest (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Elapsed":0}
{"Time":"2023-11-27T15:39:39.223392-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"PASS\n"}
{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/passpackage\t0.332s\n"}
{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Elapsed":0.333}
{"Time":"2023-11-27T15:43:53.048232-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2","Output":"=== RUN   TestPassTest2\n"}
{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2","Output":"--- PASS: TestPassTest2 (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestPassTest2","Elapsed":0}
{"Time":"2023-11-27T15:43:53.397883-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest"}
{"Time":"2023-11-27T15:43:53.397884-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest"}
{"Time":"2023-11-27T15:43:53.397885-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest"}
{"Time":"2023-11-27T15:43:53.397999-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"=== RUN   TestFailTest\n"}
{"Time":"2023-11-27T15:43:53.398176-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"    mirror_test.go:12: \n"}
{"Time":"2023-11-27T15:43:53.398181-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n"}
{"Time":"2023-11-27T15:43:53.398191-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"        \tError:      \tAn error is expected but got nil.\n"}
{"Time":"2023-11-27T15:43:53.398197-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"        \tTest:       \tTestFailTest\n"}
{"Time":"2023-11-27T15:43:53.398634-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Output":"--- FAIL: TestFailTest (0.00s)\n"}
{"Time":"2023-11-27T15:43:53.39865-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Test":"TestFailTest","Elapsed":0}
{"Time":"2023-11-27T15:43:53.398988-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Output":"FAIL\n"}
{"Time":"2023-11-27T15:43:53.399042-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/failpackage\t0.349s\n"}
{"Time":"2023-11-27T15:43:53.399052-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/failpackage","Elapsed":0.351}
`

type stdoutRedirector struct {
	originalStdout *os.File    // To keep track of the original os.Stdout
	r              *os.File    // Read end of the pipe
	w              *os.File    // Write end of the pipe
	outC           chan string // Channel to capture output
	isClosed       bool
}

func newStdoutRedirector() *stdoutRedirector {
	return &stdoutRedirector{isClosed: true}
}

func (sr *stdoutRedirector) startRedirect() {
	// Save original stdout
	sr.originalStdout = os.Stdout

	// Create a pipe for capturing output
	sr.r, sr.w, _ = os.Pipe()
	os.Stdout = sr.w

	// Initialize the output capture channel
	sr.outC = make(chan string)

	// Start capturing
	go func() {
		var buf io.ReadCloser = sr.r
		b, _ := io.ReadAll(buf)
		sr.outC <- string(b)
	}()
	sr.isClosed = false
}

func (sr *stdoutRedirector) closeRedirect() string {
	if sr.isClosed {
		return ""
	}
	// Close the write end of the pipe to finish the capture
	sr.w.Close()

	// Restore original stdout
	os.Stdout = sr.originalStdout

	// Read the captured output from the channel
	out := <-sr.outC

	sr.isClosed = true

	return out
}

// do not use with parallel tests
func captureModifierOutput(t *testing.T, input string, modifiers []TestLogModifier, config *TestLogModifierConfig, shouldCapture bool) string {
	reader := bytes.NewBufferString(input)
	output := ""
	sr := newStdoutRedirector()
	if shouldCapture {
		sr.startRedirect()
	}
	// defer close so stdout gets fixed even on test errors
	defer sr.closeRedirect()

	err := ReadAndModifyLogs(testcontext.Get(t), reader, modifiers, config)
	require.NoError(t, err, "Error reading and modifying logs")

	output = sr.closeRedirect()

	return output
}

func TestReadAndModifyLogs(t *testing.T) {
	tests := []struct {
		name     string
		config   *TestLogModifierConfig
		input    string
		expected string
	}{
		{
			name: "NonJsonTestEvents",
			config: &TestLogModifierConfig{
				IsJsonInput:            ptr.Ptr(false),
				RemoveTLogPrefix:       ptr.Ptr(false),
				ShouldImmediatelyPrint: true,
			},
			input: `non-json-test-event-line
non-json-test-event-line
non-json-test-event-line
`,
			expected: `non-json-test-event-line
non-json-test-event-line
non-json-test-event-line
`,
		},
		{
			name: "NonJsonRemovePrefix",
			config: &TestLogModifierConfig{
				IsJsonInput:            ptr.Ptr(false),
				RemoveTLogPrefix:       ptr.Ptr(true),
				ShouldImmediatelyPrint: true,
			},
			input: `non-json-test-event-line
					environment.go:1023: + ./remote-runner.test
non-json-test-event-line
`,
			expected: `non-json-test-event-line
					+ ./remote-runner.test
non-json-test-event-line
`,
		},
		{
			name: "JsonRemovePrefixStdOutput",
			config: &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       ptr.Ptr(false),
				CI:               ptr.Ptr(false),
			},

			input:    inputPassFailMixed,
			expected: "    + ./remote-runner.test\n--- PASS: TestPassTest (0.00s)\nPASS\nok  \tgithub.com/smartcontractkit/chainlink-testing-framework/passpackage\t0.332s\n--- PASS: TestPassTest2 (0.00s)\n    \n        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest\n--- FAIL: TestFailTest (0.00s)\nFAIL\nFAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/failpackage\t0.349s\n",
		},
		{
			name: "JsonRemovePrefixStdOutputOnlyErrors",
			config: &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       ptr.Ptr(true),
				CI:               ptr.Ptr(false),
			},

			input:    inputPassFailMixed,
			expected: "    \n        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest\n--- FAIL: TestFailTest (0.00s)\nFAIL\nFAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/failpackage\t0.349s\n",
		},
		{
			name: "JsonRemovePrefixStdOutputOnlyErrorsInCI",
			config: &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       ptr.Ptr(true),
				CI:               ptr.Ptr(true),
			},

			input:    inputPassFailMixed,
			expected: "FAIL  \t\x1b[31mgithub.com/smartcontractkit/chainlink-testing-framework/failpackage\x1b[0m\t0.351000\n::group::\x1b[31mTestFailTest\x1b[0m\n    \n        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest\n--- FAIL: TestFailTest (0.00s)\n::endgroup::\nFAIL\nFAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/failpackage\t0.349s\n",
		},
		{
			name: "JsonPackagePanicWithoutTestFail",
			config: &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       ptr.Ptr(true),
				CI:               ptr.Ptr(false),
			},

			input: `{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"=== RUN   TestGetImage\n"}
{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}
{"Time":"2023-11-27T15:39:39.223392-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"PASS\n"}
{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}
{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}
{"Time":"2023-11-27T15:39:38.521970916Z","Action":"start","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage","Output":"=== RUN   TestGetImage\n"}
{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestGetImage","Elapsed":0}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"pause","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager"}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"cont","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager"}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"run","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager"}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager","Output":"=== RUN   TestNewManager\n"}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager","Output":"--- PASS: TestNewManager (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.521970916Z","Action":"pass","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Test":"TestNewManager","Elapsed":0}
{"Time":"2023-11-28T11:38:06.521970916Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"PASS\n"}
{"Time":"2023-11-28T11:38:06.528992418Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"panic: Log in goroutine after TestNewManager has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n"}
{"Time":"2023-11-28T11:38:06.529024227Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\n"}
{"Time":"2023-11-28T11:38:06.529029266Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"goroutine 80 [running]:\n"}
{"Time":"2023-11-28T11:38:06.52903657Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"testing.(*common).logDepth(0xc00172c000, {0xc000052840, 0xad}, 0x3)\n"}
{"Time":"2023-11-28T11:38:06.529043322Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1022 +0x4c5\n"}
{"Time":"2023-11-28T11:38:06.529047831Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"testing.(*common).log(...)\n"}
{"Time":"2023-11-28T11:38:06.529052609Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1004\n"}
{"Time":"2023-11-28T11:38:06.529057168Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"testing.(*common).Logf(0xc00172c000, {0x18770d6?, 0x4110c5?}, {0xc0016a1190?, 0x15d9f00?, 0x1?})\n"}
{"Time":"2023-11-28T11:38:06.529061666Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1055 +0x54\n"}
{"Time":"2023-11-28T11:38:06.529068459Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"go.uber.org/zap/zaptest.testingWriter.Write({{0x7f63c4847198?, 0xc00172c000?}, 0x70?}, {0xc0017aa800?, 0xae, 0xc0016a1180?})\n"}
{"Time":"2023-11-28T11:38:06.529073568Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zaptest/logger.go:130 +0xdc\n"}
{"Time":"2023-11-28T11:38:06.529109796Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"go.uber.org/zap/zapcore.(*ioCore).Write(0xc0017808d0, {0x1, {0xc15192279f1426a5, 0x9fcc1ad, 0x2bb1d20}, {0xc001080060, 0x2c}, {0xc001080270, 0x27}, {0x1, ...}, ...}, ...)\n"}
{"Time":"2023-11-28T11:38:06.529129653Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zapcore/core.go:99 +0xb5\n"}
{"Time":"2023-11-28T11:38:06.529134612Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc0010bc820, {0x0, 0x0, 0x0})\n"}
{"Time":"2023-11-28T11:38:06.529141134Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zapcore/entry.go:253 +0x1dc\n"}
{"Time":"2023-11-28T11:38:06.529188643Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"go.uber.org/zap.(*SugaredLogger).log(0xc000246168, 0x1, {0x197ac35?, 0x19?}, {0xc0016a1140?, 0x1?, 0x1?}, {0x0, 0x0, 0x0})\n"}
{"Time":"2023-11-28T11:38:06.529202639Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/sugar.go:316 +0xec\n"}
{"Time":"2023-11-28T11:38:06.529206726Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"go.uber.org/zap.(*SugaredLogger).Warnf(...)\n"}
{"Time":"2023-11-28T11:38:06.541785942Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Output":"FAIL\tgithub.com/smartcontractkit/chainlink/v2/core/services/telemetry\t0.192s\n"}
{"Time":"2023-11-28T11:38:06.541823091Z","Action":"fail","Package":"github.com/smartcontractkit/chainlink/v2/core/services/telemetry","Elapsed":0.192}
{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage"}
{"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage","Output":"=== RUN   TestGetImage\n"}
{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}
{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Test":"TestGetImage","Elapsed":0}
{"Time":"2023-11-27T15:39:39.223392-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Output":"PASS\n"}
{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/mirror2\t0.332s\n"}
{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror2","Elapsed":0.333}`,
			expected: "--- PASS: TestGetImage (0.00s)\n--- PASS: TestNewManager (0.00s)\nPASS\npanic: Log in goroutine after TestNewManager has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n\ngoroutine 80 [running]:\ntesting.(*common).logDepth(0xc00172c000, {0xc000052840, 0xad}, 0x3)\n\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1022 +0x4c5\ntesting.(*common).log(...)\n\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1004\ntesting.(*common).Logf(0xc00172c000, {0x18770d6?, 0x4110c5?}, {0xc0016a1190?, 0x15d9f00?, 0x1?})\n\t/opt/hostedtoolcache/go/1.21.4/x64/src/testing/testing.go:1055 +0x54\ngo.uber.org/zap/zaptest.testingWriter.Write({{0x7f63c4847198?, 0xc00172c000?}, 0x70?}, {0xc0017aa800?, 0xae, 0xc0016a1180?})\n\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zaptest/logger.go:130 +0xdc\ngo.uber.org/zap/zapcore.(*ioCore).Write(0xc0017808d0, {0x1, {0xc15192279f1426a5, 0x9fcc1ad, 0x2bb1d20}, {0xc001080060, 0x2c}, {0xc001080270, 0x27}, {0x1, ...}, ...}, ...)\n\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zapcore/core.go:99 +0xb5\ngo.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc0010bc820, {0x0, 0x0, 0x0})\n\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/zapcore/entry.go:253 +0x1dc\ngo.uber.org/zap.(*SugaredLogger).log(0xc000246168, 0x1, {0x197ac35?, 0x19?}, {0xc0016a1140?, 0x1?, 0x1?}, {0x0, 0x0, 0x0})\n\t/home/runner/go/pkg/mod/go.uber.org/zap@v1.26.0/sugar.go:316 +0xec\ngo.uber.org/zap.(*SugaredLogger).Warnf(...)\nFAIL\tgithub.com/smartcontractkit/chainlink/v2/core/services/telemetry\t0.192s\n",
		},
	}
	for _, test := range tests {
		name := test.name
		config := test.config
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			require.NoError(t, config.Validate(), "Config should be valid")
			m := SetupModifiers(config)
			found := captureModifierOutput(t, input, m, config, true)
			require.Equal(t, expected, found, "Expected %v, got %v", expected, found)
		})
	}
}
