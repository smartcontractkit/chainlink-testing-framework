package exampletestpackage

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestPass(t *testing.T) {
	t.Parallel()
	t.Log("This test always passes")
}

func TestFail(t *testing.T) {
	t.Parallel()
	t.Fatal("This test always fails")
}

func TestFailLargeOutput(t *testing.T) {
	t.Parallel()
	for i := 0; i < 1000; i++ {
		t.Log("This is a log line")
	}
	t.Fatal("This test always fails")
}

func TestSubTestsAllPass(t *testing.T) {
	t.Parallel()

	t.Run("Pass1", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Run("Pass2", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})
}

func TestSubTestsAllFail(t *testing.T) {
	t.Parallel()

	t.Run("Fail1", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest always fails")
	})

	t.Run("Fail2", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest always fails")
	})
}

func TestSubTestsSomeFail(t *testing.T) {
	t.Parallel()

	t.Run("Pass", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest always fails")
	})
}

func TestSubTestsSomePanic(t *testing.T) {
	t.Parallel()

	t.Run("Pass", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Run("Panic", func(t *testing.T) {
		t.Parallel()
		panic("This subtest always panics")
	})
}

func TestFailInParentAfterSubTests(t *testing.T) {
	t.Parallel()

	t.Run("Pass1", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Run("Pass2", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Fatal("This test always fails")
}

func TestFailInParentBeforeSubTests(t *testing.T) {
	t.Parallel()

	t.Fatal("This test always fails") //nolint:revive

	t.Run("Pass1", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	t.Run("Pass2", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})
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

func TestFlakyPanic(t *testing.T) {
	t.Parallel()

	// Track if the test has run before
	stateFile := "tmp_test_flaky_panic_state"

	// If the state file does not exist, create it and fail the test
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		if err := os.WriteFile(stateFile, []byte("run once"), 0644); err != nil { //nolint:gosec
			t.Fatalf("THIS IS UNEXPECTED: failed to create state file: %v", err)
		}
		panic("This is a designed flaky test panicking as intended")
	}
	t.Cleanup(func() {
		err := os.Remove(stateFile)
		if err != nil {
			t.Fatalf("THIS IS UNEXPECTED: failed to remove state file: %v", err)
		}
	})
}

func TestRace(t *testing.T) {
	t.Parallel()
	t.Logf("This test should trigger a failure if run with the -race flag, but otherwise pass")

	var (
		numGoroutines = 100
		sharedCounter int
		wg            sync.WaitGroup
	)

	worker := func(id int) {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sharedCounter++
			_ = sharedCounter * id
		}
	}

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go worker(i)
	}

	wg.Wait()

	// Log the result
	t.Logf("Final value of sharedCounter: %d", sharedCounter)
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	deadline, ok := t.Deadline()
	if !ok {
		log.Fatal("This test should have a deadline")
	}

	t.Logf("This test will sleep %s in order to timeout", time.Until(deadline).String())
	// Sleep until the deadline
	time.Sleep(time.Until(deadline))
	t.Logf("This test should have timed out")
}

// 1) No subtests at all
func TestParentNoSubtests(t *testing.T) {
	t.Parallel()

	t.Log("No subtests, just a single test that passes.")
	// (Optional) you could also do t.Fail() or t.Fatal() to produce a fail
}

// 2) All subtests pass, no parent fail
func TestParentAllPassSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Parent does not fail, subtests all pass")

	t.Run("SubtestA", func(t *testing.T) {
		t.Parallel()
		t.Log("passes")
	})
	t.Run("SubtestB", func(t *testing.T) {
		t.Parallel()
		t.Log("passes")
	})
}

// 3) All subtests fail, no parent fail
func TestParentAllFailSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Parent does not fail, subtests all fail => typically the parent is marked fail by Go")

	t.Run("FailA", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest always fails")
	})
	t.Run("FailB", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest always fails")
	})
}

// 4) Some subtests pass, some fail, parent does NOT do its own fail
func TestParentSomeFailSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Parent does not fail, subtests partially pass/fail => parent is typically fail unless 'zeroOutParentFailsIfSubtestOnlyFails' modifies it")

	t.Run("Pass", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest passes")
	})
	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		t.Fatal("This subtest fails")
	})
}

// 5) Parent fails *after* subtests
func TestParentOwnFailAfterSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Parent fails after subtests pass => genuine parent-level failure")

	t.Run("Pass1", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})
	t.Run("Pass2", func(t *testing.T) {
		t.Parallel()
		t.Log("This subtest always passes")
	})

	// Finally, parent fails
	t.Fatal("Parent test fails after subtests pass")
}

// 6) Parent fails *before* subtests
func TestParentOwnFailBeforeSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Parent fails before subtests => subtests might not even run in real usage, or still get reported, depending on concurrency")

	t.Fatal("Parent test fails immediately")

	t.Run("WouldPassButNeverRuns", func(t *testing.T) {
		t.Parallel()
		t.Log("Normally passes, but might not even run now.")
	})
}

// 7) Nested subtests: parent -> child -> grandchild
func TestNestedSubtests(t *testing.T) {
	t.Parallel()
	t.Log("Deep nesting example")

	t.Run("Level1", func(t *testing.T) {
		t.Parallel()

		t.Run("Level2Pass", func(t *testing.T) {
			t.Parallel()
			t.Log("This sub-subtest passes")
		})

		t.Run("Level2Fail", func(t *testing.T) {
			t.Parallel()
			t.Fatal("This sub-subtest fails")
		})
	})
}

func TestParentWithFailingSubtest(t *testing.T) {
	// The parent does NOT fail. Only subtests do.
	t.Run("FailingSubtest", func(t *testing.T) {
		t.Errorf("This subtest always fails.")
	})
	t.Run("PassingSubtest", func(t *testing.T) {
		// pass
	})
}

func TestParentWithFailingParentAndSubtest(t *testing.T) {
	// Run a subtest that fails.
	t.Run("FailingSubtest", func(t *testing.T) {
		t.Errorf("This subtest always fails.")
	})
	// Run a subtest that passes.
	t.Run("PassingSubtest", func(t *testing.T) {
		// pass
	})
	// The parent test also fails.
	t.Errorf("parent fails")
}
