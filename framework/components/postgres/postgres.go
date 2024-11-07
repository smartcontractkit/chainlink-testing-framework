package postgres

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strings"
	"time"
)

const (
	User              = "chainlink"
	Password          = "thispasswordislongenough"
	Port              = "5432"
	ExposedStaticPort = "13000"
	Database          = "chainlink"
	DBVolumeName      = "postgresql_data"
)

type Input struct {
	Image      string  `toml:"image" validate:"required"`
	Port       string  `toml:"port"`
	VolumeName string  `toml:"volume_name"`
	Databases  int     `toml:"databases"`
	PullImage  bool    `toml:"pull_image"`
	Out        *Output `toml:"out"`
}

type Output struct {
	Url               string `toml:"url"`
	DockerInternalURL string `toml:"docker_internal_url"`
}

func NewPostgreSQL(in *Input) (*Output, error) {
	ctx := context.Background()

	bindPort := fmt.Sprintf("%s/tcp", Port)
	containerName := framework.DefaultTCName("postgresql")

	var sqlCommands []string
	for i := 0; i <= in.Databases; i++ {
		sqlCommands = append(sqlCommands, fmt.Sprintf("CREATE DATABASE db_%d;", i))
	}
	sqlCommands = append(sqlCommands, "ALTER USER chainlink WITH SUPERUSER;")
	initSQL := strings.Join(sqlCommands, "\n")
	initFile, err := os.CreateTemp("", "init-*.sql")
	if err != nil {
		return nil, err
	}
	if _, err := initFile.WriteString(initSQL); err != nil {
		return nil, err
	}
	if err := initFile.Close(); err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           fmt.Sprintf("%s", in.Image),
		Name:            containerName,
		Labels:          framework.DefaultTCLabels(),
		ExposedPorts:    []string{bindPort},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Env: map[string]string{
			"POSTGRES_USER":     User,
			"POSTGRES_PASSWORD": Password,
			"POSTGRES_DB":       Database,
		},
		Cmd: []string{
			"postgres", "-c", fmt.Sprintf("port=%s", Port),
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
				FileMode:          0644,
			},
		},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericVolumeMountSource{
					Name: fmt.Sprintf("%s%s", DBVolumeName, in.VolumeName),
				},
				Target: "/var/lib/postgresql/data",
			},
		},
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "127.0.0.1",
			"-U", User, "-p", Port, "-c", "select", "1", "-d", Database}).
			WithStartupTimeout(20 * time.Second).
			WithPollInterval(1 * time.Second),
	}
	var port string
	if in.Port != "" {
		port = in.Port
	} else {
		port = ExposedStaticPort
	}
	req.HostConfigModifier = func(h *container.HostConfig) {
		h.PortBindings = nat.PortMap{
			nat.Port(bindPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%s/tcp", port),
				},
			},
		}
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	return &Output{
		DockerInternalURL: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			User,
			Password,
			containerName,
			Port,
			Database,
		),
		Url: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			User,
			Password,
			host,
			ExposedStaticPort,
			Database,
		),
	}, nil
}
