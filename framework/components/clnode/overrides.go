package clnode

/*
This file contains data that need to be overridden dynamically when we setup more than one node or connect to ephemeral networks
*/

import (
	"bytes"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"os"
	"path/filepath"
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

func writeTestConfigOverrides(cfgData string) (string, error) {
	cfgPath := filepath.Join(framework.PathCLNode, "overrides.toml")
	co, err := os.Create(cfgPath)
	if err != nil {
		return "", err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return "", err
	}
	return cfgPath, nil
}

func writeUserConfigOverrides(cfgData string) (string, error) {
	cfgPath := filepath.Join(framework.PathCLNode, "user-overrides.toml")
	co, err := os.Create(cfgPath)
	if err != nil {
		return "", err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return "", err
	}
	return cfgPath, nil
}

func writeTestSecretsOverrides(cfgData string) (string, error) {
	cfgPath := filepath.Join(framework.PathCLNode, "secrets-overrides.toml")
	co, err := os.Create(cfgPath)
	if err != nil {
		return "", err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return "", err
	}
	return cfgPath, nil
}

func writeUserSecretsOverrides(cfgData string) (string, error) {
	cfgPath := filepath.Join(framework.PathCLNode, "user-secrets-overrides.toml")
	co, err := os.Create(cfgPath)
	if err != nil {
		return "", err
	}
	_, err = co.WriteString(cfgData)
	if err != nil {
		return "", err
	}
	return cfgPath, nil
}

func writeDefaultSecrets(pgOut *postgres.Output) (string, error) {
	secretsPath := filepath.Join(framework.PathCLNode, "secrets.toml")
	secretsOverrides, err := generateSecretsConfig(pgOut.Url, DefaultTestKeystorePassword)
	if err != nil {
		return "", err
	}
	sec, err := os.Create(secretsPath)
	if err != nil {
		return "", err
	}
	_, err = sec.WriteString(secretsOverrides)
	if err != nil {
		return "", err
	}
	return secretsPath, nil
}

func writeDefaultConfig(in *Input) (string, error) {
	cfgPath := filepath.Join(framework.PathCLNode, "config.toml")
	cfg, err := generateDefaultConfig(in.Node.Port)
	if err != nil {
		return "", err
	}
	overrides, err := os.Create(cfgPath)
	if err != nil {
		return "", err
	}
	_, err = overrides.WriteString(cfg)
	if err != nil {
		return "", err
	}
	return cfgPath, nil
}
