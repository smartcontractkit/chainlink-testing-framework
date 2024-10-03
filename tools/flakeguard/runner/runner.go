package runner

import (
	"fmt"
	"os"
	"os/exec"
)

type Runner struct {
	Verbose bool
	Dir     string
}

func (r *Runner) RunTests(packages []string) error {
	for _, p := range packages {
		if err := r.runSingleTest(p); err != nil {
			return fmt.Errorf("failed to run test at %s: %w", p, err)
		}
	}
	return nil
}

func (r *Runner) runSingleTest(testPackage string) error {
	fmt.Printf("Running tests for %s\n", testPackage)
	cmd := exec.Command("go", "test", testPackage)
	cmd.Dir = r.Dir

	// If you want to handle output or errors differently, you can redirect stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("test failed at %s: %w", testPackage, err)
	}
	return nil
}
