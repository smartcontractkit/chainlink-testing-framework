package test

import (
	"testing"
)

const SKIP_FLAKY_TEST = "SKIP FLAKY TEST:"

// SkipFlake adds a parseable message to the test output that indicates the test should be skipped if it fails in CI.
func SkipFlake(t *testing.T, jiraTicket, message string) {
	t.Logf("%s %s: %s\n%s", SKIP_FLAKY_TEST, t.Name(), jiraTicket, message)
}
