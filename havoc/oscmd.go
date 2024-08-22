package havoc

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
)

func ExecCmd(command string) (string, error) {
	L.Info().Interface("Command", command).Msg("Executing command")
	c := strings.Split(command, " ")
	cmd := exec.CommandContext(context.Background(), c[0], c[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()
			L.Error().
				Int("Code", exitCode).
				Msg("Command exited with status code")
			L.Error().
				Str("Out", stdout.String()).
				Str("Err", stderr.String()).
				Msg("Command output")
		}
	} else {
		L.Info().Msg("Command ran successfully")
		L.Debug().
			Str("Out", stdout.String()).
			Str("Err", stderr.String()).
			Msg("Command output")
	}
	return stdout.String(), err
}
