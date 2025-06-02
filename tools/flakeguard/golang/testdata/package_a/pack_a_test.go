package package_a

import "testing"

// This test should not get skipped by flakeguard
func TestPackAPass(t *testing.T) {
	t.Log("This is a passing test")
}

// This test should get skipped by flakeguard
func TestPackAFail(t *testing.T) {
	t.Log("This is a failing test")
	t.FailNow()
}
