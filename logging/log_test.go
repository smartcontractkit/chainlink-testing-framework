package logging

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogAfterTestEnd(t *testing.T) {
	l := GetTestLogger(t)
	go func() {
		for i := 0; i < 5000; i++ {
			l.Info().Msg("test")
		}
	}()
	time.Sleep(1 * time.Millisecond)
}

func TestGetTestLogger(t *testing.T) {
	l := GetTestLogger(t)
	l.Info().Msg("test")
	require.NotNil(t, l)
}

func TestGetTestContainersGoTestLogger(t *testing.T) {
	l := GetTestContainersGoTestLogger(t)
	require.NotNil(t, l.(CustomT).L)
}

// TestSplitStringIntoChunks tests the splitStringIntoChunks function with a string up to a million characters.
func TestSplitStringIntoChunks(t *testing.T) {
	// Create a test string with a million and 1 characters ('a').
	testString := strings.Repeat("a", 1000001)
	chunkSize := 100000
	expectedNumChunks := 11

	chunks := SplitStringIntoChunks(testString, chunkSize)

	require.Equal(t, expectedNumChunks, len(chunks), "Wrong number of chunks")

	// Check the size of each chunk.
	for i, chunk := range chunks {
		require.False(t, i < expectedNumChunks-1 && len(chunk) != chunkSize, "Chunk %d is not of expected size %d", i+1, chunkSize)

		// Check the last chunk size.
		if i == expectedNumChunks-1 {
			require.Equal(t, 1, len(chunk), "Last chunk is not of expected size")
		}
	}
}
