package flakytestpackage

import (
	"os"
	"sync"
	"testing"
)

func TestPass(t *testing.T) {
	t.Parallel()
	t.Log("This test always passes")
}

func TestFail(t *testing.T) {
	t.Parallel()
	t.Fatal("This test always fails")
}

// This test should have a 50% pass ratio
func TestFlaky(t *testing.T) {
	t.Parallel()

	// Track if the test has run before
	stateFile := "tmp_test_flaky_state"

	// If the state file does not exist, create it and fail the test
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		if err := os.WriteFile(stateFile, []byte("run once"), 0644); err != nil { //nolint:gosec
			t.Fatalf("THIS IS UNEXPECTED: failed to create state file: %v", err)
		}
		t.Fatalf("This is a designed flaky test working as intended")
	} else {
		t.Cleanup(func() {
			err := os.Remove(stateFile)
			if err != nil {
				t.Fatalf("THIS IS UNEXPECTED: failed to remove state file: %v", err)
			}
		})
	}

	t.Log("This test passes after the first run")
}

func TestSkipped(t *testing.T) {
	t.Parallel()
	t.Skip("This test is intentionally skipped")
}

func TestPanic(t *testing.T) {
	t.Parallel()
	panic("This test intentionally panics")
}

func TestRace(t *testing.T) {
	t.Parallel()
	t.Logf("This test should trigger a failure if run with the -race flag, but otherwise pass")

	var sharedCounter int
	var wg sync.WaitGroup

	// Define a worker function that accesses sharedCounter without synchronization
	worker := func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sharedCounter++
		}
	}

	// Start multiple goroutines to introduce a race condition
	const numGoroutines = 10
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go worker()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Log the result
	t.Logf("Final value of sharedCounter: %d", sharedCounter)
}
