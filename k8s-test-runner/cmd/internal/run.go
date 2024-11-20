package internal

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s-test-runner/runner"
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
	Run.Flags().BoolP("detached", "d", false, "Run in detached mode")
}

func runRunE(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	detachedMode, err := cmd.Flags().GetBool("detached")
	if err != nil {
		return err
	}

	var runnerCfg = &config.Runner{}
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

	if runnerCfg.TTLSecondsAfterFinished == 0 {
		runnerCfg.TTLSecondsAfterFinished = 600
	}

	runnerCfg.DetachedMode = detachedMode

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

func getChartOverrides(c config.Runner) map[string]interface{} {
	image := fmt.Sprintf("%s/%s:%s", c.ImageRegistryURL, c.ImageName, c.ImageTag)
	envMap := c.Envs
	if envMap == nil {
		envMap = map[string]string{}
	}
	envMap[c.TestConfigBase64EnvName] = c.TestConfigBase64

	labelsMap := c.Metadata.Labels
	if labelsMap == nil {
		labelsMap = make(map[string]string)
	}

	return map[string]interface{}{
		"namespace": c.Namespace,
		"rbac": map[string]interface{}{
			"roleName":           c.RBACRoleName,
			"serviceAccountName": c.RBACServiceAccountName,
		},
		"jobs":                    c.JobCount,
		"sync":                    c.SyncValue,
		"ttlSecondsAfterFinished": c.TTLSecondsAfterFinished,
		"test": map[string]interface{}{
			"name":    c.TestName, // Set this to your specific test name
			"timeout": c.TestTimeout,
		},
		"image":           image,
		"imagePullPolicy": "Always",
		"labels":          map[string]interface{}{},
		"annotations":     map[string]interface{}{}, // Add specific annotations if needed
		"metadata": map[string]interface{}{
			"labels": labelsMap,
		},
		"env": envMap,
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    c.ResourcesRequestsCPU,
				"memory": c.ResourcesRequestsMemory,
			},
			"limits": map[string]interface{}{
				"cpu":    c.ResourcesLimitsCPU,
				"memory": c.ResourcesLimitsMemory,
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
