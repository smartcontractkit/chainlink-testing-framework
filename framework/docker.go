package framework

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
)

func GetHost(container tc.Container) (string, error) {
	host, err := container.Host(context.Background())
	if err != nil {
		return "", err
	}
	// if localhost then force it to ipv4 localhost
	if host == "localhost" {
		host = "127.0.0.1"
	}
	return host, nil
}

func MapTheSamePort(port string) nat.PortMap {
	return nat.PortMap{
		nat.Port(port): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port,
			},
		},
	}
}

func DefaultTCLabels() map[string]string {
	return map[string]string{
		"framework": "ctf",
		"logging":   "promtail",
	}
}

func DefaultTCName(name string) string {
	return fmt.Sprintf("%s-%s", name, uuid.NewString()[0:5])
}
