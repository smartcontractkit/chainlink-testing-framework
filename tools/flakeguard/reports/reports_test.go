package reports

import (
	"bytes"
	"strings"
	"testing"
)

func TestFilterFailedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.5, Skipped: false},
		{TestName: "Test2", PassRatio: 0.9, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	failedTests := FilterFailedTests(results, 0.6)
	expected := []string{"Test1", "Test3"}

	if len(failedTests) != len(expected) {
		t.Fatalf("expected %d failed tests, got %d", len(expected), len(failedTests))
	}

	for i, test := range failedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
	}
}

func TestFilterPassedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: false},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true}, // Skipped test
	}

	passedTests := FilterPassedTests(results, 0.6)
	expected := []string{"Test1", "Test2"}

	if len(passedTests) != len(expected) {
		t.Fatalf("expected %d passed tests, got %d", len(expected), len(passedTests))
	}

	for i, test := range passedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
	}
}

func TestFilterSkippedTests(t *testing.T) {
	results := []TestResult{
		{TestName: "Test1", PassRatio: 0.7, Skipped: false},
		{TestName: "Test2", PassRatio: 1.0, Skipped: true},
		{TestName: "Test3", PassRatio: 0.3, Skipped: false},
		{TestName: "Test4", PassRatio: 0.8, Skipped: true},
	}

	skippedTests := FilterSkippedTests(results)
	expected := []string{"Test2", "Test4"}

	if len(skippedTests) != len(expected) {
		t.Fatalf("expected %d skipped tests, got %d", len(expected), len(skippedTests))
	}

	for i, test := range skippedTests {
		if test.TestName != expected[i] {
			t.Errorf("expected test %s, got %s", expected[i], test.TestName)
		}
	}
}

func TestPrintTests(t *testing.T) {
	tests := []TestResult{
		{
			TestName:    "Test1",
			TestPackage: "package1",
			PassRatio:   0.75,
			Skipped:     false,
			Runs:        4,
			Outputs:     []string{"Output1", "Output2"},
			Durations:   []float64{1.2, 0.9, 1.1, 1.0},
		},
	}

	// Use a buffer to capture the output
	var buf bytes.Buffer

	// Call PrintTests with the buffer
	PrintTests(tests, &buf)

	// Get the output as a string
	output := buf.String()
	expectedContains := []string{
		"TestName: Test1",
		"TestPackage: package1",
		"PassRatio: 0.75",
		"Skipped: false",
		"Runs: 4",
		"Durations: 1.20s, 0.90s, 1.10s, 1.00s",
		"Outputs:\nOutput1Output2",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, but it did not", expected)
		}
	}
}
