package framework

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
)

// ExecCmd executes a command and logs the output interactively
func ExecCmd(l zerolog.Logger, command string) ([]byte, error) {
	return ExecCmdWithContext(context.Background(), l, command)
}

// ExecCmdWithContext a command and logs the output interactively
func ExecCmdWithContext(ctx context.Context, l zerolog.Logger, command string) ([]byte, error) {
	return ExecCmdWithOpts(
		ctx,
		l,
		command,
		func(m string) {
			l.Debug().Str("Stream", "stdout").Msg(m)
		},
		func(m string) {
			l.Debug().Str("Stream", "stderr").Msg(m)
		},
	)
}

func ExecCmdWithOpts(ctx context.Context, l zerolog.Logger, command string, stdoutFunc func(string), stderrFunc func(string)) ([]byte, error) {
	c := strings.Split(command, " ")
	l.Info().Interface("Command", command).Msg("Executing command")
	cmd := exec.CommandContext(ctx, c[0], c[1:]...) // #nosec: G204

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// create a buffer, listen to both pipe outputs, wait them to finish and merge output
	// both log it and return merged output
	var combinedBuf strings.Builder
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		readStdPipe(stdout, func(m string) {
			stdoutFunc(m)
			combinedBuf.WriteString(m + "\n")
		})
		close(stdoutDone)
	}()
	go func() {
		readStdPipe(stderr, func(m string) {
			stderrFunc(m)
			combinedBuf.WriteString(m + "\n")
		})
		close(stderrDone)
	}()
	<-stdoutDone
	<-stderrDone

	err = cmd.Wait()
	return []byte(combinedBuf.String()), err
}

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
