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
		// 		{
		// 			name: "NonJsonTestEvents",
		// 			config: &TestLogModifierConfig{
		// 				IsJsonInput:            ptr.Ptr(false),
		// 				ShouldImmediatelyPrint: true,
		// 			},
		// 			input: `non-json-test-event-line
		// non-json-test-event-line
		// non-json-test-event-line
		// `,
		// 			expected: `non-json-test-event-line
		// non-json-test-event-line
		// non-json-test-event-line
		// `,
		// 		},
		// 		{
		// 			name: "NonJsonRemovePrefix",
		// 			config: &TestLogModifierConfig{
		// 				IsJsonInput:            ptr.Ptr(false),
		// 				RemoveTLogPrefix:       ptr.Ptr(true),
		// 				ShouldImmediatelyPrint: true,
		// 			},
		// 			input: `non-json-test-event-line
		// 			environment.go:1023: + ./remote-runner.test
		// non-json-test-event-line
		// `,
		// 			expected: `non-json-test-event-line
		// 			+ ./remote-runner.test
		// non-json-test-event-line
		// `,
		// 		},
		// 		{
		// 			name: "JsonRemovePrefix",
		// 			config: &TestLogModifierConfig{
		// 				IsJsonInput:            ptr.Ptr(true),
		// 				RemoveTLogPrefix:       ptr.Ptr(true),
		// 				OnlyErrors:             ptr.Ptr(false),
		// 				ShouldImmediatelyPrint: true,
		// 			},
		// 			input: `{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"=== RUN   TestPassTest\n"}
		// {"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"    environment.go:1023: + ./remote-runner.test\n"}
		// {"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"--- PASS: TestPassTest (0.00s)\n"}
		// {"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Elapsed":0}
		// {"Time":"2023-11-27T15:39:39.223392-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"PASS\n"}
		// {"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/passpackage\t0.332s\n"}
		// {"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Elapsed":0.333}
		// `,
		// 			expected: `{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"pause","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"cont","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223129-07:00","Action":"run","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"=== RUN   TestPassTest\n"}
		// {"Time":"2023-11-27T15:39:39.223203-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"    + ./remote-runner.test\n"}
		// {"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest","Output":"--- PASS: TestPassTest (0.00s)\n"}
		// {"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Test":"TestPassTest"}
		// {"Time":"2023-11-27T15:39:39.223392-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"PASS\n"}
		// {"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/passpackage\t0.332s\n"}
		// {"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage","Elapsed":0.333}
		// `,
		// 		},
		{
			name: "JsonRemovePrefixStdOutput",
			config: &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       ptr.Ptr(false),
				CI:               ptr.Ptr(false),
			},

			input: `{"Time":"2023-11-27T15:39:38.891004-07:00","Action":"start","Package":"github.com/smartcontractkit/chainlink-testing-framework/passpackage"}
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
`,
			expected: `=== RUN   TestPassTest
	+ ./remote-runner.test
--- PASS: TestPassTest (0.00s)
PASS
ok      github.com/smartcontractkit/chainlink-testing-framework/passpackage     0.332s`,
		},
	}
	for _, test := range tests {
		name := test.name
		config := test.config
		input := test.input
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			m := SetupModifiers(config)
			found := captureModifierOutput(t, input, m, config, true)
			require.Equal(t, expected, found, "Expected %v, got %v", input, found)
		})
	}
}

// func TestReadEvents_OnlyTestEvents(t *testing.T) {
// 	input := `{"Time":"2020-01-01T00:00:00Z","Action":"run","Package":"mypackage","Test":"TestExample","Output":"output1","Elapsed":0.1}
// {"Time":"2020-01-02T00:00:00Z","Action":"pass","Package":"mypackage","Test":"TestExample","Output":"output2","Elapsed":0.2}
// `
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 2, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 0, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }

// func TestReadEvents_MixedTestEvents(t *testing.T) {
// 	input := `{"Time":"2020-01-01T00:00:00Z","Action":"run","Package":"mypackage","Test":"TestExample","Output":"output1","Elapsed":0.1}
// non-json-test-event-line
// {"Time":"2020-01-02T00:00:00Z","Action":"pass","Package":"mypackage","Test":"TestExample","Output":"output2","Elapsed":0.2}
// `
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 2, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 1, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }

// func TestReadEvents_EmptyString(t *testing.T) {
// 	// test empty string
// 	input := ""
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 0, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 0, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)

// 	// test empty lines in string
// 	input = `

// `
// 	testEventCounter, nonTestEventCounter = lineCounterHelper(t, input)
// 	require.Exactly(t, 0, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 2, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }
