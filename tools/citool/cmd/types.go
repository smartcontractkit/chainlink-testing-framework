package cmd

type Test struct {
	Name string
	Path string
}

// CITestConf defines the configuration for running a test in a CI environment, specifying details like test ID, path, type, runner settings, command, and associated workflows.
type CITestConf struct {
	ID          string `yaml:"id" json:"id"`
	IDSanitized string `json:"id_sanitized"`
	// Name is a human-readable name for the test
	Name        string `yaml:"name" json:"name"`
	Path        string `yaml:"path" json:"path"`
	TestEnvType string `yaml:"test_env_type" json:"test_env_type"`
	// RunsOn denotes the type of GitHub actions runner to use for the test: https://docs.github.com/en/billing/managing-billing-for-github-actions/about-billing-for-github-actions#per-minute-rates
	RunsOn string `yaml:"runs_on" json:"runs_on"`
	// ChainlinkImageTypes is a list of Chainlink image variants to test with
	ChainlinkImageTypes        []string          `yaml:"chainlink_image_types" json:"chainlink_image_types"`
	TestCmd                    string            `yaml:"test_cmd" json:"test_cmd"`
	TestConfigOverrideRequired bool              `yaml:"test_config_override_required" json:"test_config_override_required"`
	TestConfigOverridePath     string            `yaml:"test_config_override_path" json:"test_config_override_path"`
	TestSecretsRequired        bool              `yaml:"test_secrets_required" json:"test_secrets_required"`
	TestEnvVars                map[string]string `yaml:"test_env_vars" json:"test_env_vars"`
	RemoteRunnerMemory         string            `yaml:"remote_runner_memory" json:"remote_runner_memory"`
	PyroscopeEnv               string            `yaml:"pyroscope_env" json:"pyroscope_env"`
	Triggers                   []string          `yaml:"triggers" json:"triggers"`
}

type Config struct {
	Tests []CITestConf `yaml:"runner-test-matrix"`
}
