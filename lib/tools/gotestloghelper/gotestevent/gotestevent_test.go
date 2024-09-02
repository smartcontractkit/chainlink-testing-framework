package gotestevent

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/clihelper"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
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
	tests := []struct {
		name     string
		event    *GoTestEvent
		expected string
	}{
		{
			name: "Error",
			event: &GoTestEvent{
				Time:    time.Now(),
				Action:  "output",
				Package: "github.com/smartcontractkit/chainlink-testing-framework/lib/failpackage",
				Test:    "TestFailTest",
				Output:  "        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest\n",
			},
			expected: "\x1b[0;31m        \tError Trace:\t/Users/blarg/git/chainlink-testing-framework/failpackage/mirror_test.go:12\n        \tError:      \tAn error is expected but got nil.\n        \tTest:       \tTestFailTest \x1b[0m\n",
		},
		{
			name: "Panic",
			event: &GoTestEvent{
				Time:    time.Now(),
				Action:  "output",
				Package: "github.com/smartcontractkit/chainlink-testing-framework/lib/failpackage",
				Test:    "TestFailTest",
				Output:  "panic: close of closed channel\n",
			},
			expected: "\x1b[0;31mpanic: close of closed channel \x1b[0m\n",
		},
		{
			name: "SegFault",
			event: &GoTestEvent{
				Time:    time.Now(),
				Action:  "output",
				Package: "github.com/smartcontractkit/chainlink-testing-framework/lib/failpackage",
				Test:    "TestFailTest",
				Output:  "[signal SIGSEGV: segmentation violation\n",
			},
			expected: "\x1b[0;31m[signal SIGSEGV: segmentation violation \x1b[0m\n",
		},
	}
	for _, test := range tests {
		name := test.name
		event := test.event
		expected := test.expected
		t.Run(name, func(t *testing.T) {
			c := &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				HidePassingTests: &clihelper.BoolFlag{IsSet: true, Value: true},
				CI:               ptr.Ptr(true),
			}

			err := HighlightErrorOutput(event, c)
			require.NoError(t, err, "Error highlighting error output")
			require.Equal(t, expected, event.Output)
		})
	}
}

func TestRemoveTestLogPrefix(t *testing.T) {
	te := &GoTestEvent{
		Time:    time.Now(),
		Action:  "output",
		Package: "github.com/smartcontractkit/chainlink-testing-framework/lib/failpackage",
		Test:    "TestFailTest",
		Output:  "    environment.go:1023: + ./remote-runner.test\n"}
	c := &TestLogModifierConfig{
		IsJsonInput:      ptr.Ptr(true),
		RemoveTLogPrefix: ptr.Ptr(true),
		HidePassingTests: &clihelper.BoolFlag{IsSet: true, Value: true},
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
				Package: "github.com/smartcontractkit/chainlink-testing-framework/lib/failpackage",
				Test:    name,
				Output:  input}
			c := &TestLogModifierConfig{
				IsJsonInput:      ptr.Ptr(true),
				RemoveTLogPrefix: ptr.Ptr(true),
				HidePassingTests: &clihelper.BoolFlag{IsSet: false, Value: false},
				CI:               ptr.Ptr(true),
			}
			require.NoError(t, c.Validate(), "Config should be valid")
			err := JsonTestOutputToStandard(te, c)
			require.NoError(t, err, "Error highlighting error output")
			require.Equal(t, 1, len(c.TestPackageMap))
			require.Equal(t, te.Package, c.TestPackageMap[te.Package].Name)
			if name == "" {
				require.Equal(t, 0, len(c.TestPackageMap[te.Package].Tests))
				require.Equal(t, 0, len(c.TestPackageMap[te.Package].TestOrder))
				require.Equal(t, false, c.TestPackageMap[te.Package].Failed)
			} else {
				require.Equal(t, 1, len(c.TestPackageMap[te.Package].Tests))
				require.Equal(t, 1, len(c.TestPackageMap[te.Package].TestOrder))
			}
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
		hidePassingTests bool
		hidePassingLogs  bool
		errorAtTopLength *int
		singlePackage    bool
	}{
		{
			name: "ShowPassingTests",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"=== RUN\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":1.02}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"PASS\n"}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected: "📦 \x1b[0;32mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage (1.02s) \x1b[0m\nabc\n::endgroup::\n",
		},
		{
			name: "AllTestsPassButWeOnlyWantToShowErrors",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"ok  \tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "",
			hidePassingTests: true,
		},
		{
			name: "ShowFailingTestsWithOnlyErrorsFalse",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- FAIL: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected: "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n",
		},
		{
			name: "ShowFailingTestsWithOnlyErrorsTrue",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- FAIL: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n",
			hidePassingTests: true,
		},
		{
			name: "HidePassingLogs",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:        "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n\x1b[0;32m✅ TestGetImage1 (0.00s) \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			hidePassingLogs: true,
		},
		{
			name: "CombinedPassFailAndCombinedOutput",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected: "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;32m✅ TestGetImage1 (0.00s) \x1b[0m\nabc\n::endgroup::\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
		},
		{
			name: "HidePassingTestsWhenBothPassAndFailExist",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			hidePassingTests: true,
		},
		{
			name: "HidePassingTestsWhenWhenBothPassAndFailExistSinglePackage",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n",
			hidePassingTests: true,
			singlePackage:    true,
		},
		{
			name: "TestPanicCompleteTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2024-05-12T00:08:06.489477083Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"panic: close of closed channel \n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦\x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror0.332s\x1b[0m::group:: \x1b[0;31m❌PANIC❌ TestGetImage2 (0.00s) \x1b[0m\nefg\npanic: close of closed channel\n::endgroup::\n",
			hidePassingTests: true,
		},
		{
			name: "TestPanicInCompleteTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2024-05-12T00:08:06.489477083Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"panic: close of closed channel \n"}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦\x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror0.332s\x1b[0m::group:: \x1b[0;31m❌PANIC❌ TestGetImage2 (0.00s) \x1b[0m\nefg\npanic: close of closed channel\n::endgroup::\n",
			hidePassingTests: true,
		},
		{
			name: "PackagePanicAfterTestPass",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- PASS: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-28T11:38:06.528992418Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"panic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected: "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n::group:: \x1b[0;32m✅ TestGetImage2 (0.00s) \x1b[0m\nefg\n::endgroup::\n\x1b[0;31mpanic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"} \x1b[0m\n",
		},
		{
			name: "PackagePanicAfterTestPassOnlyErrors",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Output":"--- PASS: TestGetImage (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- PASS: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-28T11:38:06.528992418Z","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"panic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"}\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage (0.00s) \x1b[0m\nabc\n::endgroup::\n\x1b[0;31mpanic: Log in goroutine after TestGetImage has completed: 2023-11-28T11:38:06.521Z\tWARN\tTelemetryManager.TelemetryIngressBatchClient\twsrpc@v0.7.2/uni_client.go:97\tctx error context canceled reconnecting\t{\"version\": \"2.7.0@0957729\"} \x1b[0m\n",
			hidePassingTests: true,
		},
		{
			name: "NoDropDownIfNoLogsInTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n\x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\n",
			hidePassingTests: true,
		},
		{
			name: "NoDropDownIfNoLogsInTest",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"example 1\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"example 2\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"example 3\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"    test_common.go:193: \n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"        \tError:      \tReceived unexpected error:\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"        \tTest:       \tTestWithSingleNodeEnvLocalCharts\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "📦 \x1b[0;31mgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s \x1b[0m\n::group:: \x1b[0;31m❌ TestGetImage2 (0.00s) \x1b[0m\n---❌ErrorFound❌---\n        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n        \tError:      \tReceived unexpected error:\n        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n        \tTest:       \tTestWithSingleNodeEnvLocalCharts\n---❌EndError❌---\nexample 1\nexample 2\nexample 3\n    test_common.go:193: \n        \tError Trace:\t/home/runner/work/chainlink-testing-framework/chainlink-testing-framework/k8s/e2e/common/test_common.go:193\n        \tError:      \tReceived unexpected error:\n        \t            \twaitcontainersready, no pods in 'chainlink-testing-framework-k8s-test-862b1' with selector '' after timeout '15m0s'\n        \tTest:       \tTestWithSingleNodeEnvLocalCharts\n::endgroup::\n",
			hidePassingTests: true,
			errorAtTopLength: ptr.Ptr(2),
		},
		{
			name: "SinglePackageWithSegFault",
			inputs: []string{
				`{"Time":"2024-05-12T00:07:20.069895141Z","Action":"start","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke"}`,
				`{"Time":"2024-05-12T00:07:20.32149023Z","Action":"run","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade"}`,
				`{"Time":"2024-05-12T00:07:20.321539933Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade","Output":"\u001b[90m00:07:20.32\u001b[0m \u001b[32mINF\u001b[0m Reading configs from file system\n"}`,
				`{"Time":"2024-05-12T00:07:20.32495929Z","Action":"run","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0"}`,
				`{"Time":"2024-05-12T00:08:04.238366571Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"    automation_test.go:1324: \n"}`,
				`{"Time":"2024-05-12T00:08:04.23838193Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \tError Trace:\t/home/runner/work/chainlink/chainlink/integration-tests/smoke/automation_test.go:1324\n"}`,
				`{"Time":"2024-05-12T00:08:04.23838773Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \t            \t\t\t\t/home/runner/work/chainlink/chainlink/integration-tests/smoke/automation_test.go:126\n"}`,
				`{"Time":"2024-05-12T00:08:04.23839279Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \tError:      \tReceived unexpected error:\n"}`,
				`{"Time":"2024-05-12T00:08:04.23839816Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \t            \tfailed to start CL node container err: Error response from daemon: No such image: public.ecr.aws/chainlink/chainlink:latest: failed to create container\n"}`,
				`{"Time":"2024-05-12T00:08:04.238403179Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \tTest:       \tTestAutomationNodeUpgrade/registry_2_0\n"}`,
				`{"Time":"2024-05-12T00:08:04.238408109Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"        \tMessages:   \tError deploying test environment\n"}`,
				`{"Time":"2024-05-12T00:08:06.489477083Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"panic: runtime error: invalid memory address or nil pointer dereference\n"}`,
				`{"Time":"2024-05-12T00:08:06.489485969Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"[signal SIGSEGV: segmentation violation code=0x1 addr=0xd8 pc=0x5ab986f]\n"}`,
				`{"Time":"2024-05-12T00:08:06.489491479Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"\n"}`,
				`{"Time":"2024-05-12T00:08:06.489495457Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"goroutine 1788 [running]:\n"}`,
				`{"Time":"2024-05-12T00:08:06.48950252Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env.(*ClCluster).Stop.func1()\n"}`,
				`{"Time":"2024-05-12T00:08:06.489507309Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"\t/home/runner/work/chainlink/chainlink/integration-tests/docker/test_env/cl_node_cluster.go:54 +0x2f\n"}`,
				`{"Time":"2024-05-12T00:08:06.489511628Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"golang.org/x/sync/errgroup.(*Group).Go.func1()\n"}`,
				`{"Time":"2024-05-12T00:08:06.489515956Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"\t/home/runner/go/pkg/mod/golang.org/x/sync@v0.6.0/errgroup/errgroup.go:78 +0x56\n"}`,
				`{"Time":"2024-05-12T00:08:06.489521276Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"created by golang.org/x/sync/errgroup.(*Group).Go in goroutine 354\n"}`,
				`{"Time":"2024-05-12T00:08:06.489525463Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"\t/home/runner/go/pkg/mod/golang.org/x/sync@v0.6.0/errgroup/errgroup.go:75 +0x96\n"}`,
				`{"Time":"2024-05-12T00:08:06.495576951Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Output":"FAIL\tgithub.com/smartcontractkit/chainlink/integration-tests/smoke\t46.425s\n"}`,
				`{"Time":"2024-05-12T00:08:06.495595336Z","Action":"fail","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Elapsed":46.426}`,
			},
			expected:         "::group::\x1b[0;33mIncompleteTest:TestAutomationNodeUpgrade(0.00s)\x1b[0m\x1b[90m00:07:20.32\x1b[0m\x1b[32mINF\x1b[0mReadingconfigsfromfilesystem::endgroup::::group::\x1b[0;31m❌PANIC❌TestAutomationNodeUpgrade/registry_2_0(0.00s)\x1b[0mautomation_test.go:1324:ErrorTrace:/home/runner/work/chainlink/chainlink/integration-tests/smoke/automation_test.go:1324/home/runner/work/chainlink/chainlink/integration-tests/smoke/automation_test.go:126Error:Receivedunexpectederror:failedtostartCLnodecontainererr:Errorresponsefromdaemon:Nosuchimage:public.ecr.aws/chainlink/chainlink:latest:failedtocreatecontainerTest:TestAutomationNodeUpgrade/registry_2_0Messages:Errordeployingtestenvironmentpanic:runtimeerror:invalidmemoryaddressornilpointerdereference[signalSIGSEGV:segmentationviolationcode=0x1addr=0xd8pc=0x5ab986f]goroutine1788[running]:github.com/smartcontractkit/chainlink/integration-tests/docker/test_env.(*ClCluster).Stop.func1()/home/runner/work/chainlink/chainlink/integration-tests/docker/test_env/cl_node_cluster.go:54+0x2fgolang.org/x/sync/errgroup.(*Group).Go.func1()/home/runner/go/pkg/mod/golang.org/x/sync@v0.6.0/errgroup/errgroup.go:78+0x56createdbygolang.org/x/sync/errgroup.(*Group).Goingoroutine354/home/runner/go/pkg/mod/golang.org/x/sync@v0.6.0/errgroup/errgroup.go:75+0x96::endgroup::",
			hidePassingTests: true,
			singlePackage:    true,
		},
		{
			name: "TestFailWithSIGSEGV",
			inputs: []string{
				`{"Time":"2024-05-12T00:07:20.069895141Z","Action":"start","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke"}`,
				`{"Time":"2024-05-12T00:07:20.32149023Z","Action":"run","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade"}`,
				`{"Time":"2024-05-12T00:07:20.321539933Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade","Output":"\u001b[90m00:07:20.32\u001b[0m \u001b[32mINF\u001b[0m Reading configs from file system\n"}`,
				`{"Time":"2024-05-12T00:07:20.32495929Z","Action":"run","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0"}`,
				`{"Time":"2024-05-12T00:08:06.489485969Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Output":"[signal SIGSEGV: segmentation violation code=0x1 addr=0xd8 pc=0x5ab986f]\n"}`,
				`{"Time":"2024-05-12T00:08:06.489525463Z","Action":"fail","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Test":"TestAutomationNodeUpgrade/registry_2_0","Elapsed":0}`,
				`{"Time":"2024-05-12T00:08:06.495576951Z","Action":"output","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Output":"FAIL\tgithub.com/smartcontractkit/chainlink/integration-tests/smoke\t46.425s\n"}`,
				`{"Time":"2024-05-12T00:08:06.495595336Z","Action":"fail","Package":"github.com/smartcontractkit/chainlink/integration-tests/smoke","Elapsed":46.426}`,
			},
			expected:         "📦\x1b[0;31mgithub.com/smartcontractkit/chainlink/integration-tests/smoke46.425s\x1b[0m::group::\x1b[0;33mIncompleteTest:TestAutomationNodeUpgrade(0.00s)\x1b[0m\x1b[90m00:07:20.32\x1b[0m\x1b[32mINF\x1b[0mReadingconfigsfromfilesystem::endgroup::::group::\x1b[0;31m❌PANIC❌TestAutomationNodeUpgrade/registry_2_0(0.00s)\x1b[0m[signalSIGSEGV:segmentationviolationcode=0x1addr=0xd8pc=0x5ab986f]::endgroup::",
			hidePassingTests: true,
			singlePackage:    false,
		},
		{
			name: "CombinedPassFailWithoutPanic",
			inputs: []string{
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"abc\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Output":"--- PASS: TestGetImage1 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223335-07:00","Action":"pass","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage1","Elapsed":0}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"efg\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"--- FAIL: TestGetImage2 (0.00s)\n"}`,
				`{"Time":"2023-11-27T15:39:39.223325-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Test":"TestGetImage2","Output":"nothing to see here\n"}`,
				`{"Time":"2023-11-27T15:39:39.223823-07:00","Action":"output","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Output":"FAIL\tgithub.com/smartcontractkit/chainlink-testing-framework/lib/mirror\t0.332s\n"}`,
				`{"Time":"2023-11-27T15:39:39.223871-07:00","Action":"fail","Package":"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror","Elapsed":0.333}`,
			},
			expected:         "::group::\x1b[0;33mIncompleteTest:TestGetImage2(0.00s)\x1b[0mefgnothingtoseehere::endgroup::",
			hidePassingTests: true,
			singlePackage:    true,
		},
		{
			name: "UnfinishedTestsWhenPanicOccurs",
			inputs: []string{
				`{"Time":"2024-05-16T17:29:48.996704064-05:00","Action":"start","Package":"github.com/jmank88/gotest/c"}`,
				`{"Time":"2024-05-16T17:29:48.998332751-05:00","Action":"run","Package":"github.com/jmank88/gotest/c","Test":"TestC"}`,
				`{"Time":"2024-05-16T17:29:48.998344563-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestC","Output":"=== RUN   TestC\n"}`,
				`{"Time":"2024-05-16T17:29:48.998378367-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestC","Output":"=== PAUSE TestC\n"}`,
				`{"Time":"2024-05-16T17:29:48.998383998-05:00","Action":"pause","Package":"github.com/jmank88/gotest/c","Test":"TestC"}`,
				`{"Time":"2024-05-16T17:29:48.998392303-05:00","Action":"run","Package":"github.com/jmank88/gotest/c","Test":"TestLong"}`,
				`{"Time":"2024-05-16T17:29:48.998397253-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestLong","Output":"=== RUN   TestLong\n"}`,
				`{"Time":"2024-05-16T17:29:48.998404166-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestLong","Output":"=== PAUSE TestLong\n"}`,
				`{"Time":"2024-05-16T17:29:48.998408925-05:00","Action":"pause","Package":"github.com/jmank88/gotest/c","Test":"TestLong"}`,
				`{"Time":"2024-05-16T17:29:48.998414255-05:00","Action":"cont","Package":"github.com/jmank88/gotest/c","Test":"TestC"}`,
				`{"Time":"2024-05-16T17:29:48.998418784-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestC","Output":"=== CONT  TestC\n"}`,
				`{"Time":"2024-05-16T17:29:48.998423603-05:00","Action":"cont","Package":"github.com/jmank88/gotest/c","Test":"TestLong"}`,
				`{"Time":"2024-05-16T17:29:48.998428322-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestLong","Output":"=== CONT  TestLong\n"}`,
				`{"Time":"2024-05-16T17:29:49.998588859-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestC","Output":"--- FAIL: TestC (1.00s)\n"}`,
				`{"Time":"2024-05-16T17:29:50.000739915-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Test":"TestC","Output":"panic: test panic [recovered]\n"}`,
				`{"Time":"2024-05-16T17:29:50.001589475-05:00","Action":"fail","Package":"github.com/jmank88/gotest/c","Test":"TestC","Elapsed":1}`,
				`{"Time":"2024-05-16T17:29:50.001616567-05:00","Action":"output","Package":"github.com/jmank88/gotest/c","Output":"FAIL\tgithub.com/jmank88/gotest/c\t1.005s\n"}`,
				`{"Time":"2024-05-16T17:29:50.001628149-05:00","Action":"fail","Package":"github.com/jmank88/gotest/c","Elapsed":1.005}`,
			},
			expected:         "📦\x1b[0;31mgithub.com/jmank88/gotest/c1.005s\x1b[0m::group::\x1b[0;31m❌PANIC❌TestC(1.00s)\x1b[0mpanic:testpanic[recovered]::endgroup::\x1b[0;33mIncompleteTest:TestLong(0.00s)\x1b[0m",
			hidePassingTests: true,
		},
	}

	for _, test := range tests {
		name := test.name
		expected := test.expected
		hidePassingTests := test.hidePassingTests
		hidePassingLogs := test.hidePassingLogs
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
				HidePassingTests: &clihelper.BoolFlag{IsSet: true, Value: hidePassingTests},
				HidePassingLogs:  ptr.Ptr(hidePassingLogs),
				CI:               ptr.Ptr(true),
				SinglePackage:    ptr.Ptr(singlePackage),
				ErrorAtTopLength: errorAtTopLength,
			}
			require.NoError(t, c.Validate(), "Config should be valid")
			SetupModifiers(c)
			output := genericCaptureOutput(func() {
				for _, input := range inputs {
					testEvent, err := ParseTestEvent([]byte(input))
					require.NoError(t, err, "input: ", input)
					err = JsonTestOutputToStandard(testEvent, c)
					require.NoError(t, err)
				}
			}, true)
			require.Equal(t, strings.Join(strings.Fields(expected), ""), strings.Join(strings.Fields(output), ""))
		})
	}

}
