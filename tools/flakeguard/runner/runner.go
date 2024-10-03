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

func (r *Runner) RunTests(testPaths []string) error {
	for _, path := range testPaths {
		if err := r.runSingleTest(path); err != nil {
			return fmt.Errorf("failed to run test at %s: %w", path, err)
		}
	}
	return nil
}

func (r *Runner) runSingleTest(path string) error {
	fmt.Printf("Running test for %s\n", path)
	cmd := exec.Command("go", "test", path)
	cmd.Dir = r.Dir

	// If you want to handle output or errors differently, you can redirect stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("test failed at %s: %w", path, err)
	}
	return nil
}
