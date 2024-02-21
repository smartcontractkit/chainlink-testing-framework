package gotestevent

// import (
// 	"bytes"
// 	"testing"

// 	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
// 	"github.com/stretchr/testify/require"
// )

// func lineCounterHelper(t *testing.T, input string) (int, int) {
// 	reader := bytes.NewBufferString(input)

// 	// Mock TestEventReader
// 	testEventCounter := 0
// 	the := func(te *TestEvent) error {
// 		testEventCounter++
// 		if te.Package != "mypackage" {
// 			t.Errorf("Expected 'mypackage', got '%s'", te.Package)
// 		}
// 		return nil
// 	}

// 	// Mock NonTestEventReader
// 	nonTestEventCounter := 0
// 	nteh := func(b []byte) error {
// 		nonTestEventCounter++
// 		return nil
// 	}

// 	err := ReadEvents(testcontext.Get(t), reader, the, nteh)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// 	return testEventCounter, nonTestEventCounter
// }

// func TestReadEvents_OnlyNonJsonTestEvents(t *testing.T) {
// 	input := `non-json-test-event-line
// non-json-test-event-line
// non-json-test-event-line
// `
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 0, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 3, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// 	t.Log("writing a test log")
// }

// func TestReadEvents_OnlyTestEvents(t *testing.T) {
// 	input := `{"Time":"2020-01-01T00:00:00Z","Action":"run","Package":"mypackage","Test":"TestExample","Output":"output1","Elapsed":0.1}
// {"Time":"2020-01-02T00:00:00Z","Action":"pass","Package":"mypackage","Test":"TestExample","Output":"output2","Elapsed":0.2}
// `
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 2, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 0, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }

// func TestReadEvents_MixedTestEvents(t *testing.T) {
// 	input := `{"Time":"2020-01-01T00:00:00Z","Action":"run","Package":"mypackage","Test":"TestExample","Output":"output1","Elapsed":0.1}
// non-json-test-event-line
// {"Time":"2020-01-02T00:00:00Z","Action":"pass","Package":"mypackage","Test":"TestExample","Output":"output2","Elapsed":0.2}
// `
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 2, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 1, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }

// func TestReadEvents_EmptyString(t *testing.T) {
// 	// test empty string
// 	input := ""
// 	testEventCounter, nonTestEventCounter := lineCounterHelper(t, input)
// 	require.Exactly(t, 0, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 0, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)

// 	// test empty lines in string
// 	input = `

// `
// 	testEventCounter, nonTestEventCounter = lineCounterHelper(t, input)
// 	require.Exactly(t, 0, testEventCounter, "Expected 0 test events, got %d", testEventCounter)
// 	require.Exactly(t, 2, nonTestEventCounter, "Expected 0 test events, got %d", nonTestEventCounter)
// }
