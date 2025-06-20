package cmd

import (
	"strings"
	"testing"
)

func TestFilterTestsByID(t *testing.T) {
	tests := []CITestConf{
		{ID: "run_all_in_ocr_tests_go", TestEnvType: "docker"},
		{ID: "run_all_in_ocr2_tests_go", TestEnvType: "docker"},
		{ID: "run_all_in_ocr3_tests_go", TestEnvType: "k8s_remote_runner"},
	}

	cases := []struct {
		description string
		inputIDs    string
		expectedLen int
	}{
		{"Filter by single ID", "run_all_in_ocr_tests_go", 1},
		{"Filter by multiple IDs", "run_all_in_ocr_tests_go,run_all_in_ocr2_tests_go", 2},
		{"Wildcard to include all", "*", 3},
		{"Empty ID string to include all", "", 3},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := filterTests(tests, "", "", c.inputIDs, false)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}

func TestFilterTestsIntegration(t *testing.T) {
	tests := []CITestConf{
		{ID: "run_all_in_ocr_tests_go", TestEnvType: "docker", Triggers: []string{"Run Nightly E2E Tests"}},
		{ID: "run_all_in_ocr2_tests_go", TestEnvType: "docker", Triggers: []string{"Run PR E2E Tests"}},
		{ID: "run_all_in_ocr3_tests_go", TestEnvType: "k8s_remote_runner", Triggers: []string{"Run PR E2E Tests"}},
	}

	cases := []struct {
		description   string
		inputNames    string
		inputWorkflow string
		inputTestType string
		inputIDs      string
		expectedLen   int
	}{
		{"Filter by test type and ID", "", "", "docker", "run_all_in_ocr2_tests_go", 1},
		{"Filter by trigger and test type", "", "Run PR E2E Tests", "docker", "*", 1},
		{"No filters applied", "", "", "", "*", 3},
		{"Filter mismatching all criteria", "", "Run Nightly E2E Tests", "", "", 1},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := filterTests(tests, c.inputWorkflow, c.inputTestType, c.inputIDs, false)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}

func TestRunsOnSelfHostedID(t *testing.T) {
	tests := []CITestConf{
		{ID: "foo", TestEnvType: "typeA", RunsOnSelfHosted: "runs-on/cpu=0/ram=0"},
		{ID: "bar", TestEnvType: "typeA", RunsOnSelfHosted: "runs-on/cpu=1/ram=1"},
		{ID: "baz", TestEnvType: "typeB", RunsOnSelfHosted: "runs-on=foo/cpu=2/ram=2"},
		{ID: "qux", TestEnvType: "typeC"},
	}

	cases := []struct {
		description        string
		inputTestType      string
		expectedSelfHosted bool
	}{
		{"Insert epoch time", "typeA", true},
		{"Dont insert epoch time", "typeB", true},
		{"No self-hosted label", "typeC", false},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := filterTests(tests, "", c.inputTestType, "", false)
			for _, test := range filtered {
				t.Logf("Test ID: %s, RunsOnSelfHosted: %s", test.ID, test.RunsOnSelfHosted)
				if c.expectedSelfHosted {
					if !strings.HasPrefix(test.RunsOnSelfHosted, "runs-on=") {
						t.Errorf("Expected RunsOnSelfHosted to start with 'runs-on='. Got: %s", test.RunsOnSelfHosted)
					}
				} else {
					if test.RunsOnSelfHosted != "" {
						t.Errorf("Expected RunsOnSelfHosted to be empty. Got: %s", test.RunsOnSelfHosted)
					}
				}
			}
		})
	}
}
