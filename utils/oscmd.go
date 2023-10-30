package utils

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
)

// ExecCmd executes a command and logs the output
func ExecCmd(l zerolog.Logger, command string) error {
	return ExecCmdWithContext(context.Background(), l, command)
}

// ExecCmdWithContext executes a command with ctx and logs the output
func ExecCmdWithContext(ctx context.Context, l zerolog.Logger, command string) error {
	return ExecCmdWithOptions(ctx, l, command, func(m string) {
		l.Debug().Str("Text", m).Msg("Std Pipe")
	})
}

// readStdPipe continuously read a pipe from the command
func readStdPipe(pipe io.ReadCloser, outputFunction func(string)) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if outputFunction != nil {
			outputFunction(m)
		}
	}
}

// ExecCmdWithOptions executes a command with ctx and logs the output with a custom logging func
func ExecCmdWithOptions(ctx context.Context, l zerolog.Logger, command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	l.Info().Interface("Command", c).Msg("Executing command")
	cmd := exec.CommandContext(ctx, c[0], c[1:]...)
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
