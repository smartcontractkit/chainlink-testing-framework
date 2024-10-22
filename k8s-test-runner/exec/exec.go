package exec

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

// Cmd executes os command, logging both streams
func Cmd(command string) error {
	return CmdWithStreamFunc(command, func(m string) {
		log.Info().Str("Text", m).Msg("Command output")
	})
}

// readStdPipe continuously read a pipe from the command
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

// CmdWithStreamFunc executes command with stream function
func CmdWithStreamFunc(command string, outputFunction func(string)) error {
	c := strings.Split(command, " ")
	// #nosec G204
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
