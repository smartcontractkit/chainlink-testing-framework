package package_b

import "testing"

func TestPackBFail(t *testing.T) {
	t.Log("This test should fail and get skipped by flakeguard")
	t.FailNow()
}

func TestPackBPass(t *testing.T) {
	t.Log("This test should pass and not get skipped by flakeguard")
}

func TestPackBFailTrick(t *testing.T) {
	t.Log("This test should not fail or get skipped by flakeguard")
}

func TestPackBAlreadySkipped(t *testing.T) {
	t.Skip("This test should already be skipped, and not be skipped again")
}
