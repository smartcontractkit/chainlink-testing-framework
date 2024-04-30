package gotestevent

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/clihelper"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
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
func genericCaptureOutput(fn func(), shouldCapture bool) string {
	output := ""
	sr := newStdoutRedirector()
	if shouldCapture {
		sr.startRedirect()
	}
	// defer close so stdout gets fixed even on test errors
	defer sr.closeRedirect()

	// call the function
	fn()

	output = sr.closeRedirect()

	return output
}

func TestHighlightErrorOutput(t *testing.T) {
	te := &GoTestEvent{
		Time:    time.Now(),
		Action:  "output",
		Package: "github.com/smartcontractkit/chainlink-testing-framework/failpackage",
		Test:    "TestFailTest",
		Output:  "        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest\n",
	}
	c := &TestLogModifierConfig{
		IsJsonInput:      ptr.Ptr(true),
		RemoveTLogPrefix: ptr.Ptr(true),
		OnlyErrors:       &clihelper.BoolFlag{IsSet: true, Value: true},
		CI:               ptr.Ptr(true),
	}

	err := HighlightErrorOutput(te, c)
	require.NoError(t, err, "Error highlighting error output")
	require.Equal(t, "\x1b[0;31m        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest \x1b[0m\n", te.Output)
}

func TestRemoveTestLogPrefix(t *testing.T) {
	te := &GoTestEvent{
		Time:    time.Now(),
		Action:  "output",
		Package: "github.com/smartcontractkit/chainlink-testing-framework/failpackage",
		Test:    "TestFailTest",
		Output:  "    environment.go:1023: + ./remote-runner.test\n"}
	c := &TestLogModifierConfig{
		IsJsonInput:      ptr.Ptr(true),
		RemoveTLogPrefix: ptr.Ptr(true),
		OnlyErrors:       &clihelper.BoolFlag{IsSet: true, Value: true},
		CI:               ptr.Ptr(true),
	}

	err := RemoveTestLogPrefix(te, c)
	require.NoError(t, err, "Error highlighting error output")
	require.Equal(t, "    + ./remote-runner.test\n", te.Output)
}

func TestParseOutNoise(t *testing.T) {
	tests := []struct {
		tname    string
		name     string
		input    string
		expected string
	}{
		{
			tname: "RUN",
			name:  "TestPassTest",
			input: "=== RUN   TestPassTest\n",
		},
		{
			tname: "PAUSE",
			name:  "TestPassTest",
			input: "=== PAUSE   TestPassTest\n",
		},
		{
			tname: "CONT",
			name:  "TestPassTest",
			input: "=== CONT   TestPassTest\n",
		},
		{
			tname: "PASS",
			name:  "",
			input: "PASS\n",
		},
		{
			tname: "FAIL",
			name:  "",
			input: "FAIL\n",
		},
	}
	for _, test := range tests {
		name := test.name
		input := test.input
		tname := test.tname
		t.Run(tname, func(t *testing.T) {
			te := &GoTestEvent{
				Time:    time.Now(),
				Action:  "output",
				Package: "github.com/smartcontractkit/chainlink-testing-framework/failpackage",
				Test:    name,
				Output:  input}
			c := &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       &clihelper.BoolFlag{IsSet: false, Value: false},
				CI:               ptr.Ptr(true),
			}
			require.NoError(t, c.Validate(), "Config should be valid")
			err := JsonTestOutputToStandard(te, c)
			require.NoError(t, err, "Error highlighting error output")
			require.Equal(t, 1, len(c.TestPackageMap))
			require.Equal(t, te.Package, c.TestPackageMap[te.Package].Name)
			require.Equal(t, 0, len(c.TestPackageMap[te.Package].TestLogs))
			require.Equal(t, 0, len(c.TestPackageMap[te.Package].TestOrder))
			require.Equal(t, 0, len(c.TestPackageMap[te.Package].FailedTests))
			require.Equal(t, 0, len(c.TestPackageMap[te.Package].PanicTests))
			require.Equal(t, false, c.TestPackageMap[te.Package].Failed)
			require.Equal(t, float64(0), c.TestPackageMap[te.Package].Elapsed)
			require.Equal(t, "", c.TestPackageMap[te.Package].Message)
		})
	}
}

func TestBasicPassAndFail(t *testing.T) {
	tests := []struct {
		name             string
		inputs           []string
		expected         string
		onlyErrors       bool
		errorAtTopLength *int
		singlePackage    bool
	}{
		{
			name: "ShowPassingTests",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"=== RUN\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"PASS\n"}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;32mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n",
			onlyErrors: false,
		},
		{
			name: "AllTestsPassButWeOnlyWantToShowErrors",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "",
			onlyErrors: true,
		},
		{
			name: "ShowFailingTestsWithOnlyErrorsFalse",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- FAIL: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n",
			onlyErrors: false,
		},
		{
			name: "ShowFailingTestsWithOnlyErrorsTrue",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- FAIL: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n",
			onlyErrors: true,
		},
		{
			name: "CombinedPassFailAndCombinedOutput",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage1 (0.00s) \x1b[0m\nabc\n::endgroup::\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			onlyErrors: false,
		},
		{
			name: "CombinedPassFailAndOnlyErrors",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			onlyErrors: true,
		},
		{
			name: "CombinedPassFailAndOnlyErrorsSinglePackage",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:      "::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			onlyErrors:    true,
			singlePackage: true,
		},
		{
			name: "PackagePanicAfterTestPass",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- PASS: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-28T11:38:06.528992418Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"panic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n::group:: \x1b[0;32m✅ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n\x1b[0;31mpanic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"} \x1b[0m\n",
			onlyErrors: false,
		},
		{
			name: "PackagePanicAfterTestPassOnlyErrors",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- PASS: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-28T11:38:06.528992418Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"panic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n\x1b[0;31mpanic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"} \x1b[0m\n",
			onlyErrors: true,
		},
		{
			name: "NoDropDownIfNoLogsInTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:   "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n\x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\n",
			onlyErrors: true,
		},
		{
			name: "NoDropDownIfNoLogsInTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"example 1\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"example 2\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"example 3\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"    test_common.go:193: \n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"        \tError:      \tReceived unexpected error:\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"        \tTest:       \tTestWithSingleNodeEnvLocalCharts\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\n❌ Error found:\n        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n        \tError:      \tReceived unexpected error:\n        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n        \tTest:       \tTestWithSingleNodeEnvLocalCharts\nexample 1\nexample 2\nexample 3\n    test_common.go:193: \n        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n        \tError:      \tReceived unexpected error:\n        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n        \tTest:       \tTestWithSingleNodeEnvLocalCharts\n::endgroup::\n",
			onlyErrors:       true,
			errorAtTopLength: ptr.Ptr(2),
		},
	}

	for _, test := range tests {
		name := test.name
		expected := test.expected
		onlyErrors := test.onlyErrors
		inputs := test.inputs
		errorAtTopLength := ptr.Ptr(50)
		singlePackage := test.singlePackage
		if test.errorAtTopLength != nil {
			errorAtTopLength = test.errorAtTopLength
		}
		t.Run(name, func(t *testing.T) {
			c := &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				OnlyErrors:       &clihelper.BoolFlag{IsSet: true, Value: onlyErrors},
				CI:               ptr.Ptr(true),
				SinglePackage:    ptr.Ptr(singlePackage),
				ErrorAtTopLength: errorAtTopLength,
			}
			require.NoError(t, c.Validate(), "Config should be valid")
			SetupModifiers(c)

			output := genericCaptureOutput(func() {
				for _, input := range inputs {
					testEvent, err := ParseTestEvent([]byte(input))
					require.NoError(t, err)
					err = JsonTestOutputToStandard(testEvent, c)
					require.NoError(t, err)
				}
			}, true)
			require.Equal(t, expected, output)
		})
	}
}
