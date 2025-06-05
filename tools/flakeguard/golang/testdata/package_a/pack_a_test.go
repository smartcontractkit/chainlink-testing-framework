package package_a

import (
	"testing"
)

func TestPackAPass(t *testing.T) {
	t.Log("This test should pass and not get skipped by flakeguard")
}

func TestPackAFail(t *testing.T) {
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
		t.Log("This subtest should fail and get skipped by flakeguard")
		t.FailNow()
	})

	t.Run("passing subtest", func(t *testing.T) {
		t.Log("This subtest should pass and not get skipped by flakeguard")
	})
}

func TestPackADifferentTName(a *testing.T) {
	a.Log("This tests should fail and get skipped by flakeguard")
	a.FailNow()
}
