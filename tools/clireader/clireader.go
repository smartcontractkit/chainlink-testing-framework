package clireader

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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := handler(s.Bytes()); err != nil {
				return err
			}
		}
	}
	return s.Err()
}

func DefaultReadLineHandler(b []byte) error {
	fmt.Print(string(b))
	return nil
}
