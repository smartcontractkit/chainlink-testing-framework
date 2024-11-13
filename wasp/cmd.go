package wasp

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// ExecCmd executes a command in the shell and logs its output. 
// It takes a command string as input and returns an error if the command fails to execute or if there are issues with the output streams. 
// The command's standard output and standard error are processed by a logging function that captures and logs the output in real-time. 
// If the command completes successfully, ExecCmd returns nil; otherwise, it returns an error detailing the failure.
func ExecCmd(command string) error {
	return ExecCmdWithStreamFunc(command, func(m string) {
		log.Info().Str("Text", m).Msg("Command output")
	})
}

// readStdPipe reads lines from the provided io.ReadCloser pipe and passes each line to the specified streamFunc. 
// If streamFunc is nil, the lines will be ignored. The function continues reading until the pipe is closed or an error occurs. 
// It is typically used to handle output from command execution in a streaming manner.
func readStdPipe(pipe io.ReadCloser, streamFunc func(string)) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		if streamFunc != nil {
			streamFunc(m)
		}
	}
}

// ExecCmdWithStreamFunc executes a command specified by the command string and streams its output 
// to the provided outputFunction. The outputFunction is called with each line of output from both 
// standard output and standard error streams. The function returns an error if there is an issue 
// starting the command, creating pipes, or waiting for the command to complete. 
// If the command executes successfully, it will return nil.
func ExecCmdWithStreamFunc(command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	cmd := exec.Command(c[0], c[1:]...)
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
