package client

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func ExecCmd(command string) error {
	return ExecCmdWithContext(context.Background(), command)
}

func ExecCmdWithContext(ctx context.Context, command string) error {
	return ExecCmdWithOptions(ctx, command, func(m string) {
		log.Debug().Str("Text", m).Msg("Std Pipe")
	})
}

// readStdPipe continuously reads from a given pipe (either stdout or stderr)
// and processes the output line by line using the provided outputFunction.
// It handles lines of any length dynamically without the need for a large predefined buffer.
func readStdPipe(pipe io.ReadCloser, outputFunction func(string)) {
	reader := bufio.NewReader(pipe)
	var output []rune

	for {
		// ReadLine tries to return a single line, not including the end-of-line bytes.
		// The returned line may be incomplete if the line's too long for the buffer.
		// isPrefix will be true if the line is longer than the buffer.
		chunk, isPrefix, err := reader.ReadLine()

		// Handle any errors that occurred during the read.
		if err != nil {
			// Log any error that's not an EOF (end of file).
			if err != io.EOF {
				log.Warn().Err(err).Msg("Error while reading standard pipe, this can be caused by really long logs and can be ignored if nothing else is wrong.")
			}
			break
		}

		// Append the chunk to the output buffer.
		// bytes.Runes converts the byte slice to a slice of runes, handling multi-byte characters.
		output = append(output, bytes.Runes(chunk)...)

		// If isPrefix is false, we've reached the end of the line and can process it.
		if !isPrefix {
			// Call the output function with the complete line if it's defined.
			if outputFunction != nil {
				outputFunction(string(output))
			}
			// Reset output to an empty slice for reading the next line.
			output = output[:0]
		}
	}
}

func ExecCmdWithOptions(ctx context.Context, command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	cmd := exec.CommandContext(ctx, c[0], c[1:]...) // #nosec: G204
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go readStdPipe(stderr, outputFunction)
	go readStdPipe(stdout, outputFunction)
	return cmd.Wait()
}
