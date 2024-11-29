package reports

import "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/codeowners"

// MapTestResultsToOwners maps test results to their code owners based on the TestPath and CODEOWNERS file.
func MapTestResultsToOwners(report *TestReport, codeOwnersPath string) error {
	// Parse the CODEOWNERS file
	codeOwnerPatterns, err := codeowners.Parse(codeOwnersPath)
	if err != nil {
		return err
	}

	// Assign owners to each test result
	for i, result := range report.Results {
		if result.TestPath != "NOT FOUND" {
			report.Results[i].CodeOwners = codeowners.FindOwners(result.TestPath, codeOwnerPatterns)
		} else {
			// Mark owners as unknown for unmapped tests
			report.Results[i].CodeOwners = []string{"UNKNOWN"}
		}
	}

	return nil
}
