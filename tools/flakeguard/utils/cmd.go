package utils

import (
	"bytes"
	"os/exec"
)

type CmdOutput struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func ExecuteCmd(name string, args ...string) (*CmdOutput, error) {
	cmd := exec.Command(name, args...)
	out := &CmdOutput{}
	cmd.Stdout = &out.Stdout
	cmd.Stderr = &out.Stderr
	err := cmd.Run()
	return out, err
}
