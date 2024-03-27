package clihelper

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// CliOutputHandler is a function that handles a line of output
type CliOutputHandler func([]byte) error

// Readline reads from the reader and calls the handler for each line
func ReadLine(ctx context.Context, reader io.Reader, handler CliOutputHandler) error {
	s := bufio.NewScanner(reader)
	for s.Scan() {
		// Check for context cancellation at the start of each loop iteration.
		// This allows the function to respond promptly to cancellation requests.
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Pass the scanned line to the handler function.
		// If the handler encounters an error, return the error to the caller.
		if err := handler(s.Bytes()); err != nil {
			return err
		}
	}
	return s.Err()
}

func DefaultReadLineHandler(b []byte) error {
	fmt.Print(string(b))
	return nil
}
