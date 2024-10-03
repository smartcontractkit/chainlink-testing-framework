package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	Verbose  bool   // If true, provides detailed logging.
	Dir      string // Directory to run commands in.
	Count    int    // Number of times to run the tests.
	UseRace  bool   // Enable race detector.
	FailFast bool   // Stop on first test failure.
}

// RunTests executes the tests for each provided package.
func (r *Runner) RunTests(packages []string) error {
	for _, p := range packages {
		if err := r.runSingleTest(p); err != nil {
			return fmt.Errorf("failed to run test at %s: %w", p, err)
		}
	}
	return nil
}

// runSingleTest executes the test command for a single test package.
func (r *Runner) runSingleTest(testPackage string) error {
	args := []string{"test"}
	if r.Count > 0 {
		args = append(args, "-count", fmt.Sprint(r.Count))
	}
	if r.UseRace {
		args = append(args, "-race")
	}
	if r.FailFast {
		args = append(args, "-failfast")
	}
	args = append(args, testPackage)

	// Construct the command and display it
	cmd := exec.Command("go", args...)
	cmd.Dir = r.Dir
	cmdString := fmt.Sprintf("Executing command: go %s", strings.Join(args, " "))
	fmt.Println(cmdString)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("test failed at %s: %w", testPackage, err)
	}
	return nil
}
