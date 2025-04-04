package framework

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Config struct {
	LokiURL               string
	LokiTenantID          string
	LokiBasicAuthUsername string
	LokiBasicAuthPassword string
}

// The promtailConfig function to substitute values and write the final file
func promtailConfig() (string, error) {
	// Define the configuration as a string
	configTemplate := `
server:
  http_listen_port: 9080
  grpc_listen_port: 0

clients:
  - url: "{{ .LokiURL }}"
    tenant_id: "{{ .LokiTenantID }}"
    basic_auth:
      username: "{{ .LokiBasicAuthUsername }}"
      password: "{{ .LokiBasicAuthPassword }}"

positions:
  filename: /tmp/positions.yaml

scrape_configs:
  - job_name: flog_scrape
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
        filters:
          - name: label
            values: ["logging=promtail"]
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/(.*)'
        target_label: 'container'
      - target_label: job
        replacement: "ctf"
`

	lokiURL := os.Getenv("LOKI_URL")
	lokiTenantID := os.Getenv("LOKI_TENANT_ID")

	if lokiURL == "" {
		lokiURL = "http://host.docker.internal:3030/loki/api/v1/push"
	}
	if lokiTenantID == "" {
		lokiTenantID = "promtail"
	}

	lokiBasicAuth := os.Getenv("LOKI_BASIC_AUTH")
	var lokiBasicAuthUsername string
	var lokiBasicAuthPassword string
	if lokiBasicAuth != "" {
		authParts := strings.SplitN(lokiBasicAuth, ":", 2)
		if len(authParts) != 2 {
			return "", errors.New("LOKI_BASIC_AUTH must be in the format 'user:password'")
		}
		lokiBasicAuthUsername = authParts[0]
		lokiBasicAuthPassword = authParts[1]
	}

	secrets := Config{
		LokiURL:               lokiURL,
		LokiTenantID:          lokiTenantID,
		LokiBasicAuthUsername: lokiBasicAuthUsername,
		LokiBasicAuthPassword: lokiBasicAuthPassword,
	}

	filePath := PathRoot + "/promtail-config.yml"

	// Create the file where the promtail config will be written
	configFile, err := os.CreateTemp("", "promtail-config.yml")
	if err != nil {
		return "", fmt.Errorf("could not create promtail-config.yml file: %w", err)
	}
	defer configFile.Close()

	tmpl, err := template.New("promtail").Parse(configTemplate)
	if err != nil {
		return "", fmt.Errorf("could not parse promtail config template: %w", err)
	}

	err = tmpl.Execute(configFile, secrets)
	if err != nil {
		return "", fmt.Errorf("could not execute promtail config template: %w", err)
	}
	L.Debug().Str("Path", filePath).Msg("Promtail configuration written to")

	return configFile.Name(), nil
}

func NewPromtail() error {
	// since this container is dynamic we write it in Go, but it is a part of observability stack
	// hence, we never remove it with TESTCONTAINERS_RYUK_DISABLED but only with "ctf obs d"
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	pcn, err := promtailConfig()
	if err != nil {
		return err
	}

	cmd := make([]string, 0)
	cmd = append(cmd, "-config.file=/etc/promtail/promtail-config.yml")
	if os.Getenv("CTF_PROMTAIL_DEBUG") != "" {
		cmd = append(cmd, "-log.level=debug")
	}

	req := testcontainers.ContainerRequest{
		Image:        "grafana/promtail:latest",
		ExposedPorts: []string{"9080/tcp"},
		Name:         "promtail",
		Cmd:          cmd,
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      pcn,
				ContainerFilePath: "/etc/promtail/promtail-config.yml",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForHTTP("/ready").WithPort("9080").WithStartupTimeout(5 * time.Minute),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.Mounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/lib/docker/containers",
					Target: "/var/lib/docker/containers",
				},
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
			}
		},
	}

	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Reuse:            true,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("could not start Promtail container: %w", err)
	}
	return nil
}
