package test_test

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/test"
)

func TestSkipFlake(t *testing.T) {
	test.SkipFlake(t, "EXAMPLE-123", "This test is flaky")
	t.Fail()
}

// fail test on purpose if FLAKE_FORCE is set
func TestSkipFailOnPurposeWithEnvVar(t *testing.T) {
	if os.Getenv("FORCE_FAILURE") != "" {
		t.Fail()
	}
}
