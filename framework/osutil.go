package framework

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// ExecCmd executes a command and logs the output interactively
func ExecCmd(command string) ([]byte, error) {
	return ExecCmdWithContext(context.Background(), command)
}

// ExecCmdWithContext a command and logs the output interactively
func ExecCmdWithContext(ctx context.Context, command string) ([]byte, error) {
	return ExecCmdWithOpts(
		ctx,
		command,
		func(m string) {
			L.Debug().Str("Stream", "stdout").Msg(m)
		},
		func(m string) {
			L.Debug().Str("Stream", "stderr").Msg(m)
		},
	)
}

func ExecCmdWithOpts(ctx context.Context, command string, stdoutFunc func(string), stderrFunc func(string)) ([]byte, error) {
	c := strings.Split(command, " ")
	L.Info().Interface("Command", command).Msg("Executing command")
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
	combinedBufMu := &sync.Mutex{}
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		readStdPipe(stdout, func(m string) {
			stdoutFunc(m)
			combinedBufMu.Lock()
			combinedBuf.WriteString(m + "\n")
			combinedBufMu.Unlock()
		})
		close(stdoutDone)
	}()
	go func() {
		readStdPipe(stderr, func(m string) {
			stderrFunc(m)
			combinedBufMu.Lock()
			combinedBuf.WriteString(m + "\n")
			combinedBufMu.Unlock()
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
