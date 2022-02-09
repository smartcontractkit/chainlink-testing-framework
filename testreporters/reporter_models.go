// package testreporters holds all the tools necessary to report on tests that are run utilizing the testsetups package
package testreporters

// TestReporter
type TestReporter interface {
	WriteReport(folderPath string) error
}
