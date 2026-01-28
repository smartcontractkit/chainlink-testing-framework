package framework_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

func runCmd(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	t.Logf("Executing: %s %v in %s", name, args, dir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output:\n%s", string(output))
		t.Fatalf("Command failed: %v", err)
	}
	t.Logf("Command output:\n%s", string(output))
	return string(output)
}

type TestCfg struct {
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
}

// TestSmokeGenerateDevEnv top-down approach tests until all the environment variations aren't stable
func TestSmokeGenerateDevEnv(t *testing.T) {
	tests := []struct {
		name        string
		cliName     string
		productName string
		outputDir   string
		nodes       int
	}{
		// all variations of env, test they can be generated and create valid services
		// that can pass healthchecks
		{
			name:        "basic 2 nodes env",
			cliName:     "tcli",
			productName: "myproduct",
			outputDir:   "test-env",
			nodes:       2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				runCmd(t, tt.outputDir, tt.cliName, `down`)
				// remove this line if test fails and see what can't be compiled
				os.RemoveAll(tt.outputDir)
			})
			cg, err := framework.NewEnvBuilder(
				tt.cliName,
				tt.nodes,
				tt.productName,
			).
				OutputDir(tt.outputDir).
				Build()
			require.NoError(t, err)
			err = cg.Write()
			require.NoError(t, err)
			err = cg.WriteServices()
			require.NoError(t, err)
			err = cg.WriteFakes()
			require.NoError(t, err)
			err = cg.WriteProducts()
			require.NoError(t, err)
			runCmd(t, filepath.Join(tt.outputDir, "cmd", tt.cliName), `go`, `install`, `.`)
			runCmd(t, tt.outputDir, tt.cliName, `up`)

			data, err := os.ReadFile(filepath.Join(tt.outputDir, "env-out.toml"))
			require.NoError(t, err)

			decoder := toml.NewDecoder(strings.NewReader(string(data)))

			var cfg TestCfg
			err = decoder.Decode(&cfg)
			require.NoError(t, err)

			clClients, err := clclient.New(cfg.NodeSets[0].Out.CLNodes)
			require.NoError(t, err)

			pollEvery := 2 * time.Second
			timeout := 2 * time.Minute

			// verify that all the services are healthy
			require.Eventually(t, func() bool {
				failedServices := false
				for _, c := range clClients {
					r, _, err := c.Health()
					require.NoError(t, err)
					t.Log(r)
					for _, d := range r.Data {
						if d.Attributes.Status != "passing" {
							t.Logf("CL service %s is not healthy: %s", d.Attributes.Name, d.Attributes.Output)
							failedServices = true
						}
					}
				}
				return !failedServices
			}, timeout, pollEvery)
		})
	}
}
