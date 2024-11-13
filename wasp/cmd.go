package wasp

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// ExecCmd executes the given command string in a shell environment.
// It logs the command's output using the default logging mechanism.
// It returns any error encountered during the execution of the command.
func ExecCmd(command string) error {
	return ExecCmdWithStreamFunc(command, func(m string) {
		log.Info().Str("Text", m).Msg("Command output")
	})
}

// readStdPipe reads lines from the provided io.ReadCloser pipe and passes each line to the streamFunc callback.
// If streamFunc is nil, the lines are not processed. This function is typically used to handle output streams
// from executed commands, allowing real-time processing of command output.
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

// ExecCmdWithStreamFunc executes a shell command and streams its output to the provided outputFunction.
// The command's standard output and standard error are captured and passed to outputFunction line by line.
// It returns any error encountered during the execution of the command.
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
