package test_env

import (
	tc "github.com/testcontainers/testcontainers-go"
)

type EnvComponent struct {
	ContainerName    string   `json:"containerName"`
	ContainerImage   string   `json:"containerImage"`
	ContainerVersion string   `json:"containerVersion"`
	Networks         []string `json:"networks"`
	Container        tc.Container
}

type EnvComponentOption = func(c *EnvComponent)

func WithContainerName(name string) EnvComponentOption {
	return func(c *EnvComponent) {
		if name != "" {
			c.ContainerName = name
		}
	}
}
