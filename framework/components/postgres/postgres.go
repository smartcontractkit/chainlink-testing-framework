package postgres

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strings"
	"time"
)

type Input struct {
	Image     string  `toml:"image" validate:"required" default:"postgres"`
	Tag       string  `toml:"tag" validate:"required" default:"15.6"`
	PullImage bool    `toml:"pull_image" default:"true"`
	User      string  `toml:"user" validate:"required" default:"chainlink"`
	Password  string  `toml:"password" validate:"required" default:"thispasswordislongenough"`
	Database  string  `toml:"database" validate:"required" default:"chainlink"`
	Databases int     `toml:"databases" validate:"required" default:"20"`
	Port      string  `toml:"port" validate:"required" default:"5432"`
	Out       *Output `toml:"out"`
}

type Output struct {
	Url               string `toml:"url"`
	DockerInternalURL string `toml:"docker_internal_url"`
}

func NewPostgreSQL(in *Input) (*Output, error) {
	ctx := context.Background()

	bindPort := fmt.Sprintf("%s/tcp", in.Port)

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
		Image:           fmt.Sprintf("%s:%s", in.Image, in.Tag),
		Name:            containerName,
		Labels:          framework.DefaultTCLabels(),
		ExposedPorts:    []string{bindPort},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Env: map[string]string{
			"POSTGRES_USER":     in.User,
			"POSTGRES_PASSWORD": in.Password,
			"POSTGRES_DB":       in.Database,
		},
		Cmd: []string{
			"postgres", "-c", fmt.Sprintf("port=%s", in.Port),
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      initFile.Name(),
				ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
				FileMode:          0644,
			},
		},
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "127.0.0.1",
			"-U", in.User, "-p", in.Port, "-c", "select", "1", "-d", in.Database}).
			WithStartupTimeout(20 * time.Second),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}
	return &Output{
		DockerInternalURL: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			in.User,
			in.Password,
			containerName,
			in.Port,
			in.Database,
		),
		Url: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			in.User,
			in.Password,
			host,
			mp.Port(),
			in.Database,
		),
	}, nil
}
