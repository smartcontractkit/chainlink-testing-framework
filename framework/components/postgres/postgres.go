package postgres

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type Input struct {
	User     string  `toml:"user" validate:"required"`
	Password string  `toml:"password" validate:"required"`
	Database string  `toml:"database" validate:"required"`
	Port     string  `toml:"port" validate:"required"`
	Out      *Output `toml:"out"`
}

type Output struct {
	ContainerName string
	Url           string `toml:"url"`
}

func NewPostgreSQL(in *Input) (*Output, error) {
	ctx := context.Background()

	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	containerName := framework.DefaultTCName("postgresql")

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.6",
		Name:         containerName,
		Labels:       framework.DefaultTCLabels(),
		ExposedPorts: []string{bindPort},
		Networks:     []string{framework.DefaultNetworkName},
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
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "127.0.0.1",
			"-U", in.User, "-p", in.Port, "-c", "select", "1", "-d", in.Database}).
			WithStartupTimeout(10 * time.Second),
	}
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	//host, err := framework.GetHost(c)
	//if err != nil {
	//	return nil, err
	//}
	//mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	//if err != nil {
	//	return nil, err
	//}
	return &Output{
		Url: fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			in.User,
			in.Password,
			containerName,
			in.Port,
			in.Database,
		),
	}, nil
}
