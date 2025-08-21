package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

// reportOptions influence how reports are aggregated together
type reportOptions struct {
	maxPassRatio         float64
	reportID             string
	rerunOfReportID      string
	projectPath          string
	goProjectName        string
	raceDetection        bool
	excludedTests        []string
	selectedTests        []string
	jsonOutputPaths      []string
	branchName           string
	baseSha              string
	headSha              string
	repoURL              string
	codeownersPath       string
	gitHubWorkflowName   string
	gitHubWorkflowRunURL string
}

// TestReportOption is a functional option for configuring the aggregation process.
type TestReportOption func(*reportOptions)

func WithProjectPath(projectPath string) TestReportOption {
	return func(opts *reportOptions) {
		opts.projectPath = projectPath
	}
}

func WithGoProject(goProject string) TestReportOption {
	return func(opts *reportOptions) {
		opts.goProjectName = goProject
	}
}

func WithReportID(reportID string) TestReportOption {
	return func(opts *reportOptions) {
		opts.reportID = reportID
	}
}

func WithRerunOfReportID(rerunOfReportID string) TestReportOption {
	return func(opts *reportOptions) {
		opts.rerunOfReportID = rerunOfReportID
	}
}

func WithGeneratedReportID(genReportID bool) TestReportOption {
	return func(opts *reportOptions) {
		if !genReportID {
			return
		}
		uuid, err := uuid.NewRandom()
		if err != nil {
			panic(fmt.Errorf("error generating random report id: %w", err))
		}
		opts.reportID = uuid.String()
	}
}

func WithBranchName(branchName string) TestReportOption {
	return func(opts *reportOptions) {
		opts.branchName = branchName
	}
}

// WithHeadSha sets the head SHA for the aggregated report.
func WithHeadSha(headSha string) TestReportOption {
	return func(opts *reportOptions) {
		opts.headSha = headSha
	}
}

// WithBaseSha sets the base SHA for the aggregated report.
func WithBaseSha(baseSha string) TestReportOption {
	return func(opts *reportOptions) {
		opts.baseSha = baseSha
	}
}

// WithRepoURL sets the repository URL for the aggregated report.
func WithRepoURL(repoURL string) TestReportOption {
	return func(opts *reportOptions) {
		opts.repoURL = repoURL
	}
}

func WithRepoPath(repoPath string) TestReportOption {
	return func(opts *reportOptions) {
		opts.repoURL = repoPath
	}
}

func WithCodeOwnersPath(codeOwnersPath string) TestReportOption {
	return func(opts *reportOptions) {
		opts.codeownersPath = codeOwnersPath
	}
}

// WithGitHubWorkflowName sets the GitHub workflow name for the aggregated report.
func WithGitHubWorkflowName(githubWorkflowName string) TestReportOption {
	return func(opts *reportOptions) {
		opts.gitHubWorkflowName = githubWorkflowName
	}
}

// WithGitHubWorkflowRunURL sets the GitHub workflow run URL for the aggregated report.
func WithGitHubWorkflowRunURL(githubWorkflowRunURL string) TestReportOption {
	return func(opts *reportOptions) {
		opts.gitHubWorkflowRunURL = githubWorkflowRunURL
	}
}

// WithMaxPassRatio sets the maximum pass ratio for the aggregated report.
func WithMaxPassRatio(maxPassRatio float64) TestReportOption {
	return func(opts *reportOptions) {
		opts.maxPassRatio = maxPassRatio
	}
}

func WithGoRaceDetection(raceDetection bool) TestReportOption {
	return func(opts *reportOptions) {
		opts.raceDetection = raceDetection
	}
}

func WithExcludedTests(excludedTests []string) TestReportOption {
	return func(opts *reportOptions) {
		opts.excludedTests = excludedTests
	}
}

func WithSelectedTests(selectedTests []string) TestReportOption {
	return func(opts *reportOptions) {
		opts.selectedTests = selectedTests
	}
}

func WithJSONOutputPaths(jsonOutputPaths []string) TestReportOption {
	return func(opts *reportOptions) {
		opts.jsonOutputPaths = jsonOutputPaths
	}
}

// NewTestReport creates a new TestReport based on the provided test results and optional aggregate settings.
func NewTestReport(results []TestResult, opts ...TestReportOption) (TestReport, error) {
	defaultOpts := &reportOptions{
		reportID:             "",
		branchName:           "",
		headSha:              "",
		baseSha:              "",
		repoURL:              "",
		codeownersPath:       "",
		gitHubWorkflowName:   "",
		gitHubWorkflowRunURL: "",
		maxPassRatio:         1.0, // default value
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	r := TestReport{
		ID:                   defaultOpts.reportID,
		RerunOfReportID:      defaultOpts.rerunOfReportID,
		ProjectPath:          defaultOpts.projectPath,
		GoProject:            defaultOpts.goProjectName,
		BranchName:           defaultOpts.branchName,
		HeadSHA:              defaultOpts.headSha,
		BaseSHA:              defaultOpts.baseSha,
		RepoURL:              defaultOpts.repoURL,
		GitHubWorkflowName:   defaultOpts.gitHubWorkflowName,
		GitHubWorkflowRunURL: defaultOpts.gitHubWorkflowRunURL,
		RaceDetection:        defaultOpts.raceDetection,
		ExcludedTests:        defaultOpts.excludedTests,
		SelectedTests:        defaultOpts.selectedTests,
		MaxPassRatio:         defaultOpts.maxPassRatio,
		JSONOutputPaths:      defaultOpts.jsonOutputPaths,
		Results:              results,
	}

	r.GenerateSummaryData()

	// Map test results to paths
	err := MapTestResultsToPaths(&r, r.ProjectPath)
	if err != nil {
		return r, fmt.Errorf("error mapping test results to paths: %w", err)
	}

	// Map test results to code owners if a codeowners file is provided
	if defaultOpts.codeownersPath != "" {
		err = MapTestResultsToOwners(&r, defaultOpts.codeownersPath)
		if err != nil {
			return r, fmt.Errorf("error mapping test results to code owners: %w", err)
		}
	}

	// Set the report ID for each test result
	for i := range r.Results {
		r.Results[i].ReportID = r.ID
	}

	return r, nil
}

// TestReport reports on the parameters and results of one to many test runs
type TestReport struct {
	ID                   string       `json:"id"`
	RerunOfReportID      string       `json:"rerun_of_report_id"` // references the ID of the original/base report from which this re-run was created.
	ProjectPath          string       `json:"project_path"`
	GoProject            string       `json:"go_project"`
	BranchName           string       `json:"branch_name,omitempty"`
	HeadSHA              string       `json:"head_sha,omitempty"`
	BaseSHA              string       `json:"base_sha,omitempty"`
	RepoURL              string       `json:"repo_url,omitempty"`
	GitHubWorkflowName   string       `json:"github_workflow_name,omitempty"`
	GitHubWorkflowRunURL string       `json:"github_workflow_run_url,omitempty"`
	SummaryData          *SummaryData `json:"summary_data"`
	RaceDetection        bool         `json:"race_detection"`
	ExcludedTests        []string     `json:"excluded_tests,omitempty"`
	SelectedTests        []string     `json:"selected_tests,omitempty"`
	Results              []TestResult `json:"results,omitempty"`
	FailedLogsURL        string       `json:"failed_logs_url,omitempty"`
	JSONOutputPaths      []string     `json:"-"` // go test -json outputs from runs
	// MaxPassRatio is the maximum flakiness ratio allowed for a test to be considered not flaky
	MaxPassRatio float64 `json:"max_pass_ratio,omitempty"`
}

// SaveToFile saves the test report to a JSON file at the given path.
// It returns an error if there's any issue with marshaling the report or writing to the file.
func (testReport *TestReport) SaveToFile(outputPath string) error {
	// Create directory path if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	jsonData, err := json.MarshalIndent(testReport, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling test results to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0600); err != nil {
		return fmt.Errorf("error writing test results to file: %w", err)
	}

	return nil
}

func (tr *TestReport) PrintGotestsumOutput(w io.Writer, format string) error {
	if len(tr.JSONOutputPaths) == 0 {
		fmt.Fprintf(w, "No JSON test output paths found in test report\n")
		return nil
	}

	for _, path := range tr.JSONOutputPaths {
		cmdStr := fmt.Sprintf("cat %q | gotestsum --raw-command --format %q -- cat", path, format)
		cmd := exec.Command("bash", "-c", cmdStr)

		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("gotestsum command failed for file %s: %w\nOutput: %s", path, err, outBuf.String())
		}

		fmt.Fprint(w, outBuf.String())
		fmt.Fprint(w, "\n------------------------------------------\n\n")
	}
	return nil
}

// GenerateSummaryData generates a summary of a report's test results
func (testReport *TestReport) GenerateSummaryData() {
	var runs, testRunCount, passes, fails, skips, panickedTests, racedTests, flakyTests int

	// Map to hold unique test names that were entirely skipped
	uniqueSkippedTestsMap := make(map[string]struct{})

	for _, result := range testReport.Results {
		runs += result.Runs
		if result.Runs > testRunCount {
			testRunCount = result.Runs
		}
		passes += result.Successes
		fails += result.Failures
		skips += result.Skips

		if result.Runs == 0 && result.Skipped {
			uniqueSkippedTestsMap[result.TestName] = struct{}{}
		}

		if result.Panic {
			panickedTests++
			flakyTests++
		} else if result.Race {
			racedTests++
			flakyTests++
		} else if !result.Skipped && result.Runs > 0 && result.PassRatio < testReport.MaxPassRatio {
			flakyTests++
		}
	}

	// Calculate the unique count of skipped tests
	uniqueSkippedTestCount := len(uniqueSkippedTestsMap)

	// Calculate the raw pass ratio
	passRatio := passRatio(passes, runs)

	// Calculate the raw flake ratio
	totalTests := len(testReport.Results)
	flakeRatio := flakeRatio(flakyTests, totalTests)

	passRatioStr := formatRatio(passRatio)
	flakeTestRatioStr := formatRatio(flakeRatio)

	testReport.SummaryData = &SummaryData{
		UniqueTestsRun:         totalTests,
		UniqueSkippedTestCount: uniqueSkippedTestCount,
		TestRunCount:           testRunCount,
		PanickedTests:          panickedTests,
		RacedTests:             racedTests,
		FlakyTests:             flakyTests,
		FlakyTestPercent:       flakeTestRatioStr,

		TotalRuns:   runs,
		PassedRuns:  passes,
		FailedRuns:  fails,
		SkippedRuns: skips,
		PassPercent: passRatioStr,
	}
}
