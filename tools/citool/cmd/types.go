package cmd

type Test struct {
	Name string
	Path string
}

// CITestConf defines the configuration for running a test in a CI environment, specifying details like test ID, path, type, runner settings, command, and associated workflows.
type CITestConf struct {
	ID                         string            `yaml:"id" json:"id"`
	IDSanitized                string            `json:"id_sanitized"`
	Path                       string            `yaml:"path" json:"path"`
	TestEnvType                string            `yaml:"test_env_type" json:"test_env_type"`
	RunsOn                     string            `yaml:"runs_on" json:"runs_on"`
	TestCmd                    string            `yaml:"test_cmd" json:"test_cmd"`
	TestConfigOverrideRequired bool              `yaml:"test_config_override_required" json:"test_config_override_required"`
	TestConfigOverridePath     string            `yaml:"test_config_override_path" json:"test_config_override_path"`
	TestSecretsRequired        bool              `yaml:"test_secrets_required" json:"test_secrets_required"`
	TestEnvVars                map[string]string `yaml:"test_env_vars" json:"test_env_vars"`
	RemoteRunnerMemory         string            `yaml:"remote_runner_memory" json:"remote_runner_memory"`
	PyroscopeEnv               string            `yaml:"pyroscope_env" json:"pyroscope_env"`
	Workflows                  []string          `yaml:"workflows" json:"workflows"`
}

type Config struct {
	Tests []CITestConf `yaml:"runner-test-matrix"`
}
