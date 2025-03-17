package framework

import (
	"os"
)

// HostDockerInternal returns host.docker.internal that works both locally and in GHA
func HostDockerInternal() string {
	if os.Getenv(EnvVarCI) == "true" {
		return "http://172.17.0.1"
	}
	return "http://host.docker.internal"
}
