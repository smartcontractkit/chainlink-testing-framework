package internal

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/runner"
	"github.com/spf13/cobra"
)

var Run = &cobra.Command{
	Use:   "run",
	RunE:  runRunE,
	Short: "Run",
}

func init() {
	Run.PersistentFlags().StringP(
		"config",
		"c",
		"",
		"Path to TOML config",
	)
}

func runRunE(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	var runnerCfg = &config.RemoteRunner{}
	err = config.Read(configPath, "", runnerCfg)
	if err != nil {
		return err
	}

	var testCfgBase64 string
	if runnerCfg.TestConfigFilePath != "" {
		testCfgBase64, err = fileAsBase64(runnerCfg.TestConfigFilePath)
		if err != nil {
			return err
		}
		runnerCfg.TestConfigBase64 = testCfgBase64
	}

	if runnerCfg.SyncValue == "" {
		runnerCfg.SyncValue = fmt.Sprintf("a%s", uuid.NewString()[0:5])
	}

	p, err := runner.NewK8sTestRun(runnerCfg, getChartOverrides(*runnerCfg))
	if err != nil {
		return errors.Wrapf(err, "error creating test in k8s")
	}

	err = p.Run()
	if err != nil {
		return errors.Wrapf(err, "error running test in k8s")
	}

	return nil
}

func getChartOverrides(c config.RemoteRunner) map[string]interface{} {
	image := fmt.Sprintf("%s/%s:%s", c.ImageRegistryURL, c.ImageName, c.ImageTag)
	envMap := c.Envs
	if envMap == nil {
		envMap = map[string]string{}
	}
	envMap[c.TestConfigBase64EnvName] = c.TestConfigBase64

	return map[string]interface{}{
		"namespace": c.Namespace,
		"jobs":      c.JobCount,
		"sync":      c.SyncValue,
		"test": map[string]interface{}{
			"name":    c.TestName, // Set this to your specific test name
			"timeout": c.TestTimeout,
		},
		"image":           image,
		"imagePullPolicy": "Always",
		"labels": map[string]interface{}{
			"app": "wasp",
		},
		"annotations": map[string]interface{}{}, // Add specific annotations if needed
		"env":         envMap,
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "512Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "512Mi",
			},
		},
		"nodeSelector": map[string]interface{}{}, // Specify node selector if needed
		"tolerations":  []interface{}{},          // Specify tolerations if needed
		"affinity":     map[string]interface{}{}, // Specify affinity rules if needed
	}
}

func fileAsBase64(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "error reading file %s", filePath)
	}
	return base64.StdEncoding.EncodeToString(content), nil
}
