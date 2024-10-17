package clnode

/*
This file contains data that need to be overridden dynamically when we setup more than one node or connect to ephemeral networks
*/

import (
	"bytes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"os"
)

import (
	"text/template"
)

type OverridesConfig struct {
	HTTPPort      string
	SecureCookies bool
}

const defaultConfigTmpl = `
[Log]
Level = 'info'

[WebServer]
HTTPWriteTimeout = '30s'
SecureCookies = false
HTTPPort = {{.HTTPPort}}

[WebServer.TLS]
HTTPSPort = 0

[JobPipeline]
[JobPipeline.HTTPRequest]
DefaultTimeout = '10s'
`

func generateDefaultConfig(port string) (string, error) {
	config := OverridesConfig{
		HTTPPort:      port,
		SecureCookies: false,
	}
	tmpl, err := template.New("toml").Parse(defaultConfigTmpl)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, config)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

func writeTestConfigOverrides(cfgData string) (*os.File, error) {
	co, err := os.CreateTemp("", "overrides.toml")
	if err != nil {
		return nil, err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return nil, err
	}
	return co, nil
}

func writeUserConfigOverrides(cfgData string) (*os.File, error) {
	co, err := os.CreateTemp("", "user-overrides.toml")
	if err != nil {
		return nil, err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return nil, err
	}
	return co, nil
}

func writeTestSecretsOverrides(cfgData string) (*os.File, error) {
	co, err := os.CreateTemp("", "secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return nil, err
	}
	return co, nil
}

func writeUserSecretsOverrides(cfgData string) (*os.File, error) {
	co, err := os.CreateTemp("", "user-secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return nil, err
	}
	return co, nil
}

func writeDefaultSecrets(pgOut *postgres.Output) (*os.File, error) {
	secretsOverrides, err := generateSecretsConfig(pgOut.DockerInternalURL, DefaultTestKeystorePassword)
	if err != nil {
		return nil, err
	}
	sec, err := os.CreateTemp("", "secrets.toml")
	if err != nil {
		return nil, err
	}
	_, err = sec.WriteString(secretsOverrides)
	if err != nil {
		return nil, err
	}
	return sec, nil
}

func writeDefaultConfig(in *Input) (*os.File, error) {
	cfg, err := generateDefaultConfig(in.Node.Port)
	if err != nil {
		return nil, err
	}
	overrides, err := os.CreateTemp("", "config.toml")
	if err != nil {
		return nil, err
	}
	_, err = overrides.WriteString(cfg)
	if err != nil {
		return nil, err
	}
	return overrides, nil
}
