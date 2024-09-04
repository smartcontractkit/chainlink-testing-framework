package clihelper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

func lineCounterHelper(t *testing.T, input string) int {
	reader := bytes.NewBufferString(input)

	// Mock line handler
	lineCounter := 0
	lc := func(b []byte) error {
		lineCounter++
		return nil
	}

	err := ReadLine(testcontext.Get(t), reader, lc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	return lineCounter
}

func TestReadLine(t *testing.T) {
	input := `non-json-test-event-line
non-json-test-event-line
non-json-test-event-line
`
	lineCounter := lineCounterHelper(t, input)
	require.Exactly(t, 3, lineCounter, "Expected 3 lines, got %d", lineCounter)
}
