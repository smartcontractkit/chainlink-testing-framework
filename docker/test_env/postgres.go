package test_env

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type PostgresDb struct {
	EnvComponent
	User         string
	Password     string
	DbName       string
	Port         string
	ImageVersion string
	l            zerolog.Logger
	t            *testing.T
}

type PostgresDbOption = func(c *PostgresDb)

// Sets custom container name if name is not empty
func WithPostgresDbContainerName(name string) PostgresDbOption {
	return func(c *PostgresDb) {
		if name != "" {
			c.ContainerName = name
		}
	}
}

func WithPostgresImageVersion(version string) PostgresDbOption {
	return func(c *PostgresDb) {
		if version != "" {
			c.ImageVersion = version
		}
	}
}

func NewPostgresDb(networks []string, opts ...PostgresDbOption) *PostgresDb {
	pg := &PostgresDb{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "postgres-db", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		User:         "postgres",
		Password:     "mysecretpassword",
		DbName:       "testdb",
		Port:         "5432",
		ImageVersion: "15.3",
		l:            log.Logger,
	}
	for _, opt := range opts {
		opt(pg)
	}
	return pg
}

func (pg *PostgresDb) WithTestLogger(t *testing.T) *PostgresDb {
	pg.l = logging.GetTestLogger(t)
	pg.t = t
	return pg
}

func (pg *PostgresDb) StartContainer() error {
	req := pg.getContainerRequest()
	l := tc.Logger
	if pg.t != nil {
		l = logging.CustomT{
			T: pg.t,
			L: pg.l,
		}
	}
	c, err := tc.GenericContainer(context.Background(), tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
		Reuse:            true,
		Logger:           l,
	})
	if err != nil {
		return err
	}
	pg.Container = c

	pg.l.Info().Str("containerName", pg.ContainerName).
		Msg("Started Postgres DB container")

	return nil
}

func (pg *PostgresDb) getContainerRequest() *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         pg.ContainerName,
		Image:        fmt.Sprintf("postgres:%s", pg.ImageVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", pg.Port)},
		Env: map[string]string{
			"POSTGRES_USER":     pg.User,
			"POSTGRES_DB":       pg.DbName,
			"POSTGRES_PASSWORD": pg.Password,
		},
		Networks: pg.Networks,
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "localhost",
			"-U", pg.User, "-c", "select", "1", "-d", pg.DbName}).
			WithStartupTimeout(10 * time.Second),
	}
}
