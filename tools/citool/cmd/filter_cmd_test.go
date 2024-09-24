package cmd

import (
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
			filtered := filterTests(tests, "", "", c.inputIDs, "", false)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}

func TestFilterTestsByChainlinkImageType(t *testing.T) {
	tests := []CITestConf{
		{ChainlinkImageTypes: []string{"arm64"}, TestEnvType: "docker"},
		{ChainlinkImageTypes: []string{"amd64"}, TestEnvType: "docker"},
		{ChainlinkImageTypes: []string{"plugins", "arm64"}, TestEnvType: "k8s_remote_runner"},
		{ChainlinkImageTypes: []string{}, TestEnvType: "docker"}, // Empty image type should default to amd64
	}

	cases := []struct {
		description              string
		inputChainlinkImageTypes string
		expectedLen              int
	}{
		{"Filter by single type", "arm64", 2},
		{"Filter by single type once", "amd64", 2},
		{"Filter by non-existent type", "nonsense", 0},
		{"Empty type string to include all", "", 4},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := filterTests(tests, "", "", "", c.inputChainlinkImageTypes, false)
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
		inputTrigger  string
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
			filtered := filterTests(tests, c.inputTrigger, c.inputTestType, c.inputIDs, "", false)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}
