package internal

import (
	"encoding/base64"
	"fmt"
	"os"

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
	}

	if runnerCfg.TestConfigBase64 != "" {
		testCfgBase64 = runnerCfg.TestConfigBase64
	}

	testConfigEnvName := fmt.Sprintf("test.%s", runnerCfg.TestConfigEnvName)

	image := fmt.Sprintf("%s/%s:%s", runnerCfg.ImageRegistryURL, runnerCfg.ImageName, runnerCfg.ImageTag)

	p, err := runner.NewK8sTestRun(&runner.Config{
		Namespace: runnerCfg.Namespace,
		KeepJobs:  runnerCfg.KeepJobs,
		HelmValues: map[string]string{
			"image":                      image,
			"test.name":                  runnerCfg.TestName,
			"test.timeout":               runnerCfg.TestTimeout,
			"env.wasp.log_level":         runnerCfg.WaspLogLevel,
			"jobs":                       runnerCfg.WaspJobs,
			testConfigEnvName:            testCfgBase64,
			"env.TEST_LOG_LEVEL":         "debug",
			"env.MERCURY_TEST_LOG_LEVEL": "debug",
		},
		ChartPath: runnerCfg.ChartPath,
	})
	if err != nil {
		return errors.Wrapf(err, "error creating test in k8s")
	}

	err = p.Run()
	if err != nil {
		return errors.Wrapf(err, "error running test in k8s")
	}

	return nil
}

func fileAsBase64(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "error reading file %s", filePath)
	}
	return base64.StdEncoding.EncodeToString(content), nil
}
