package client

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadStdPipeWithLongString(t *testing.T) {
	// Create a string with a million characters ('a').
	longString := strings.Repeat("a", 1000000)

	// Use an io.Pipe to simulate the stdout or stderr pipe.
	reader, writer := io.Pipe()

	// Channel to communicate errors from the writing goroutine.
	errChan := make(chan error, 1)

	// Write the long string to the pipe in a goroutine.
	go func() {
		_, err := writer.Write([]byte(longString))
		if err != nil {
			// Send any errors to the main test goroutine via the channel.
			errChan <- err
		}
		writer.Close()
		errChan <- nil // Send nil to indicate successful write.
	}()

	// Variable to store the output from the readStdPipe function.
	var output string
	outputFunction := func(s string) {
		output = s
	}

	// Call the readStdPipe function with the reader part of the pipe.
	readStdPipe(reader, outputFunction)

	// Check for errors from the write goroutine.
	err := <-errChan
	require.NoError(t, err, "Failed to write to pipe")
	require.Equal(t, longString, output, "Output did not match the input long string")
}
