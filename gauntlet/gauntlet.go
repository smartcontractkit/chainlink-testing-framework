// Package gauntlet enables the framework to interface with the chainlink gauntlet project
package gauntlet

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	execDir string
)

// Gauntlet contains helpful data to run gauntlet commands
type Gauntlet struct {
	exec          string
	Network       string
	NetworkConfig map[string]string
}

// NewGauntlet Sets up a gauntlet struct and checks if the yarn executable exists.
func NewGauntlet() (*Gauntlet, error) {
	yarn, err := exec.LookPath("yarn")
	if err != nil {
		return &Gauntlet{}, errors.New("'yarn' is not installed")
	}
	log.Debug().Str("PATH", yarn).Msg("Executable Path")
	os.Setenv("SKIP_PROMPTS", "true")
	g := &Gauntlet{
		exec:          yarn,
		NetworkConfig: make(map[string]string),
	}
	g.GenerateRandomNetwork()
	return g, nil
}

// Flag returns a string formatted in the expected gauntlet's flag form
func (g *Gauntlet) Flag(flag, value string) string {
	return fmt.Sprintf("--%s=%s", flag, value)
}

func (g *Gauntlet) SetWorkingDir(wrkDir string) {
	execDir = wrkDir
}

// GenerateRandomNetwork Creates and sets a random network prepended with test
func (g *Gauntlet) GenerateRandomNetwork() {
	r := uuid.NewString()[0:8]
	t := time.Now().UnixMilli()
	g.Network = fmt.Sprintf("test%v%s", t, r)
	log.Debug().Str("Network", g.Network).Msg("Generated Network Name")
}

type ExecCommandOptions struct {
	ErrHandling       []string
	CheckErrorsInRead bool
	RetryCount        int
	RetryDelay        time.Duration
}

// ExecCommand Executes a gauntlet command with the provided arguments.
//
//	It will also check for any errors you specify in the output via the errHandling slice.
func (g *Gauntlet) ExecCommand(args []string, options ExecCommandOptions) (string, error) {
	output := ""
	// append gauntlet and network to args since it is always needed
	updatedArgs := append([]string{"gauntlet"}, args...)
	updatedArgs = insertArg(updatedArgs, 2, g.Flag("network", g.Network))
	printArgs(updatedArgs)

	cmd := exec.Command(g.exec, updatedArgs...) // #nosec G204
	if execDir != "" {
		cmd.Dir = execDir
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return output, err
	}

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		log.Info().Str("stdout", line).Msg("Gauntlet")
		output = fmt.Sprintf("%s%s", output, line)
		if options.CheckErrorsInRead {
			rerr := checkForErrors(options.ErrHandling, output)
			if rerr != nil {
				return output, rerr
			}
		}
		line, err = reader.ReadString('\n')
	}

	reader = bufio.NewReader(stderr)
	line, err = reader.ReadString('\n')
	for err == nil {
		log.Info().Str("stderr", line).Msg("Gauntlet")
		output = fmt.Sprintf("%s%s", output, line)
		if options.CheckErrorsInRead {
			rerr := checkForErrors(options.ErrHandling, output)
			if rerr != nil {
				return output, rerr
			}
		}
		line, err = reader.ReadString('\n')
	}

	rerr := checkForErrors(options.ErrHandling, output)
	if rerr != nil {
		return output, rerr
	}

	if strings.Compare("EOF", err.Error()) > 0 {
		return output, err
	}

	// catch any exit codes
	err = cmd.Wait()

	log.Debug().Msg("Gauntlet command Completed")
	return output, err
}

// ExecCommandWithRetries Some commands are safe to retry and in ci this can be even more so needed.
func (g *Gauntlet) ExecCommandWithRetries(args []string, options ExecCommandOptions) (string, error) {
	var output string
	var err error
	if options.RetryDelay == 0 {
		// default to 5 seconds
		options.RetryDelay = time.Second * 5
	}
	err = retry.Do(
		func() error {
			output, err = g.ExecCommand(args, options)
			return err
		},
		retry.Delay(options.RetryDelay),
		retry.MaxDelay(options.RetryDelay),
		retry.Attempts(uint(options.RetryCount)),
	)

	return output, err
}

// WriteNetworkConfigMap write a network config file for gauntlet testing.
func (g *Gauntlet) WriteNetworkConfigMap(networkDirPath string) error {
	file := filepath.Join(networkDirPath, fmt.Sprintf(".env.%s", g.Network))
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range g.NetworkConfig {
		log.Debug().Str(k, v).Msg("Gauntlet .env config value:")
		_, err = f.WriteString(fmt.Sprintf("\n%s=%s", k, v))
		if err != nil {
			return err
		}
	}
	return nil
}

// checkForErrors Loops through provided err slice to see if the error exists in the output.
func checkForErrors(errHandling []string, line string) error {
	for _, e := range errHandling {
		if strings.Contains(line, e) {
			log.Debug().Str("Error", line).Msg("Gauntlet Error Found")
			return fmt.Errorf("found a gauntlet error")
		}
	}
	return nil
}

// insertArg inserts an argument into the args slice
func insertArg(args []string, index int, valueToInsert string) []string {
	if len(args) <= index { // nil or empty slice or after last element
		return append(args, valueToInsert)
	}
	args = append(args[:index+1], args[index:]...) // index < len(a)
	args[index] = valueToInsert
	return args
}

// printArgs prints all the gauntlet args being used in a call to gauntlet
func printArgs(args []string) {
	out := "yarn"
	for _, arg := range args {
		out = fmt.Sprintf("%s %s", out, arg)

	}
	log.Info().Str("Command", out).Msg("Gauntlet")
}

func (g *Gauntlet) AddNetworkConfigVar(k string, v string) {
	g.NetworkConfig[k] = v
}
