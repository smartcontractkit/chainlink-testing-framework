package package_a

import "testing"

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

func TestPackAAlreadySkipped(t *testing.T) {
	t.Skip("This test should already be skipped, and not be skipped again")
}
