package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/test"
)

type outputConfig struct {
	excludeFlakes bool
	failedTestLog string
}

// TestOutput is the struct corresponding to the JSON output from `go test -json`
type TestOutput struct {
	Time    string   `json:"Time,omitempty"`
	Action  string   `json:"Action,omitempty"`
	Package string   `json:"Package,omitempty"`
	Test    string   `json:"Test,omitempty"`
	Output  *string  `json:"Output,omitempty"`
	Elapsed *float64 `json:"Elapsed,omitempty"`
}

func main() {
	config := handleFlags()
	readInput(config)
}

// handleFlags parses the flags and returns a config struct
func handleFlags() outputConfig {
	// Define flags
	excludeFlakesPtr := flag.Bool("excludeFlakes", false, "Exclude any test failures that aren't marked with flake")
	var failedTestLogPtr string

	// Custom flag set to parse only if excludeFlakes is set
	flagSet := flag.NewFlagSet("failedTestLog", flag.ExitOnError)
	flagSet.StringVar(&failedTestLogPtr, "failedTestLog", "", "File name to log failed tests")

	// Parse flags
	flag.Parse()

	// Additional parsing if excludeFlakes is set
	if *excludeFlakesPtr {
		err := flagSet.Parse(flag.Args())
		if err != nil {
			fmt.Println("Error parsing flag:", err)
			os.Exit(1)
		}
	}

	return outputConfig{
		excludeFlakes: *excludeFlakesPtr,
		failedTestLog: failedTestLogPtr,
	}
}

// readInput reads the JSON output from `go test -json` and outputs modified JSON
func readInput(c outputConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var testOutput TestOutput
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &testOutput); err != nil {
			// when we can't decode the json we just print it back out as is and continue
			os.Stderr.WriteString(line)
			continue
		}

		// Do functions that read the line here
		if c.excludeFlakes {
			excludeFlakesRead(testOutput)
		}

		// Re-encode as JSON and output
		modifiedJSON, err := json.Marshal(testOutput)
		if err != nil {
			// Handle JSON encoding error
			os.Stderr.WriteString(fmt.Sprintf("JSON encode error: %e\n", err))
			continue
		}
		_, err = os.Stdout.Write(modifiedJSON)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error writing line: %e\nOriginal line: %s\n", err, line))
		}
		os.Stdout.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Reading standard input: %e", err))
	}

	excludeFlakesExit(c)
}

// FailedTest Capture failed tests data
type FailedTest struct {
	Package string
	Name    string
	Flake   bool
	Pass    bool
}

var failedTestMap = map[string]*FailedTest{}

func flakeKey(t TestOutput) string {
	return fmt.Sprintf("%s/%s", t.Package, t.Test)
}

// excludeFlakesRead checks the log line output if the test is marked as flaky and if it failed
func excludeFlakesRead(t TestOutput) {
	key := flakeKey(t)
	if t.Output != nil {
		if strings.Contains(*t.Output, test.SKIP_FLAKY_TEST) {
			_, exists := failedTestMap[key]
			if exists {
				failedTestMap[key].Flake = true
			} else {
				failedTestMap[key] = &FailedTest{
					Package: t.Package,
					Name:    t.Test,
					Flake:   true,
					Pass:    true,
				}
			}
		}
	}
	// test is a failure, mark it as such
	if t.Action == "fail" && t.Test != "" {
		// Check if key exists
		_, exists := failedTestMap[key]
		if exists {
			failedTestMap[key].Pass = false
		} else {
			failedTestMap[key] = &FailedTest{
				Package: t.Package,
				Name:    t.Test,
				Flake:   false,
				Pass:    false,
			}
		}
	}
}

// excludeFlakesExit checks if any tests failed that aren't marked as flaky and exits with an error if so
func excludeFlakesExit(c outputConfig) {
	failedNonFlakes := []string{}
	for _, v := range failedTestMap {
		if !v.Flake && !v.Pass {
			failedNonFlakes = append(failedNonFlakes, fmt.Sprintf("%s/%s", v.Package, v.Name))
		}
	}
	if len(failedNonFlakes) > 0 {
		logFailedTests(c, failedNonFlakes)
		os.Exit(1)
	}
}

// logFailedTests logs failed tests to a file if a filename was provided in the flags
func logFailedTests(c outputConfig, failedNonFlakes []string) {
	if c.failedTestLog != "" {
		file, err := os.Create(c.failedTestLog)
		if err != nil {
			fmt.Println("Error creating file:", err)
			os.Exit(1)
		}
		defer file.Close()

		for _, line := range failedNonFlakes {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(1)
			}
		}
	}
}
