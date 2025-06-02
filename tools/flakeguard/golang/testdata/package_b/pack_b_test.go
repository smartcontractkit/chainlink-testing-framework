package package_b

import "testing"

// This test should get skipped by flakeguard
func TestPackBFail(t *testing.T) {
	t.Log("This is a failing test")
	t.FailNow()
}

// This test should not get skipped by flakeguard
func TestPackBPass(t *testing.T) {
	t.Log("This is a passing test")
}
