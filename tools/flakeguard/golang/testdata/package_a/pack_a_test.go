package package_a

import (
	"fmt"
	"testing"
)

func TestPackAPass(t *testing.T) {
	t.Log("This test should pass and not get skipped by flakeguard")
}

func TestPackAFail(t *testing.T) {
	// SHOULD BE SKIPPED
	t.Log("This test should fail and get skipped by flakeguard")
	t.FailNow()
}

func TestPackAFailTrick(t *testing.T) {
	t.Log("This test should not fail or get skipped by flakeguard")
}

func TestPackASkippedAlready(t *testing.T) {
	t.Skip("This test is already skipped, and should not be skipped again")
}

func TestPackAFailSubTest(t *testing.T) {
	t.Run("failing subtest", func(t *testing.T) {
		// SHOULD BE SKIPPED
		t.Log("This subtest should fail and get skipped by flakeguard")
		t.FailNow()
	})

	t.Run("passing subtest", func(t *testing.T) {
		t.Log("This subtest should pass and not get skipped by flakeguard")
	})
}

func TestPackAFailSubTestDynamicName(t *testing.T) {
	for i := range 3 {
		t.Run(fmt.Sprintf("subtest %d", i), func(t *testing.T) {
			if i == 1 {
				// SHOULD BE SKIPPED
				t.Logf("This subtest should fail and get skipped by flakeguard, iteration %d", i)
				t.FailNow()
			}
			t.Logf("This subtest should pass and not get skipped by flakeguard, iteration %d", i)
		})
	}
}

func TestPackAFailNestedSubTests(t *testing.T) {
	t.Run("parent subtest", func(t *testing.T) {
		t.Run("child subtest", func(t *testing.T) {
			// SHOULD BE SKIPPED
			t.Log("This child subtest should fail and get skipped by flakeguard")
			t.FailNow()
		})
	})
}

func TestPackAFailHelperFunction(t *testing.T) {
	helperFunctionWithSubTests(t)
}

func helperFunctionWithSubTests(t *testing.T) {
	t.Run("parent subtest", func(t *testing.T) {
		t.Run("failing child subtest", func(t *testing.T) {
			// SHOULD BE SKIPPED
			t.Log("This first child subtest should fail and get skipped by flakeguard")
			t.FailNow()
		})

		t.Run("passing child subtest", func(t *testing.T) {
			t.Log("This second child subtest should pass and not get skipped by flakeguard")
		})
	})
}
