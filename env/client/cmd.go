package client

import (
	"bufio"
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

func ExecCmdWithOptions(ctx context.Context, command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
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
