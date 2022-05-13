package soak_runner

//revive:disable:dot-imports
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/gomega"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Builds the go tests to run, and returns a path to it, along with remote config options
func buildGoTests() string {
	exePath := filepath.Join(utils.ProjectRoot, "remote.test")
	compileCmd := exec.Command("go", "test", "-c", utils.SoakRoot, "-o", exePath) // #nosec G204
	compileCmd.Env = os.Environ()
	compileCmd.Env = append(compileCmd.Env, "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	log.Info().Str("Test Directory", utils.SuiteRoot).Msg("Compiling tests")
	compileOut, err := compileCmd.Output()
	log.Debug().
		Str("Output", string(compileOut)).
		Str("Command", compileCmd.String()).
		Msg("Ran command")
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Env: %s\nCommand: %s\nCommand Output: %s", compileCmd.Env, compileCmd.String(), compileOut))

	_, err = os.Stat(exePath)
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Expected '%s' to exist", exePath))
	return exePath
}
