package cmd

type Test struct {
	Name string
	Path string
}

// CITestConf defines the configuration for running a test in a CI environment, specifying details like test ID, path, type, runner settings, command, and associated workflows.
type CITestConf struct {
	ID                         string            `yaml:"id" json:"id"`
	IDSanitized                string            `json:"idSanitized"`
	Path                       string            `yaml:"path" json:"path"`
	TestEnvType                string            `yaml:"test_env_type" json:"test_env_type"`
	RunsOn                     string            `yaml:"runs-on" json:"runsOn"`
	TestCmd                    string            `yaml:"test-cmd" json:"testCmd"`
	TestConfigOverrideRequired bool              `yaml:"test-config-override-required" json:"testConfigOverrideRequired"`
	TestSecretsRequired        bool              `yaml:"test-secrets-required" json:"testSecretsRequired"`
	TestInputs                 map[string]string `yaml:"test_inputs" json:"test_inputs"`
	RemoteRunnerMemory         string            `yaml:"remote-runner-memory" json:"remoteRunnerMemory"`
	PyroscopeEnv               string            `yaml:"pyroscope-env" json:"pyroscopeEnv"`
	Workflows                  []string          `yaml:"workflows" json:"workflows"`
}

type Config struct {
	Tests []CITestConf `yaml:"runner-test-matrix"`
}
