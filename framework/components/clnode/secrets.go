package clnode

import (
	"bytes"
	"text/template"
)

type Config struct {
	DatabaseURL string
	Keystore    string
}

const dbTmpl = `[Database]
URL = '{{.DatabaseURL}}'

[Password]
Keystore = '{{.Keystore}}'
`

func generateSecretsConfig(connString, password string) (string, error) {
	// Create the configuration with example values
	config := Config{
		DatabaseURL: connString,
		Keystore:    password,
	}
	tmpl, err := template.New("toml").Parse(dbTmpl)
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
