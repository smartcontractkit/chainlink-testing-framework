package test_env

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	InternalPort string
	ExternalPort string
	InternalURL  *url.URL
	ExternalURL  *url.URL
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

func WithPostgresDbName(name string) PostgresDbOption {
	return func(c *PostgresDb) {
		if name != "" {
			c.DbName = name
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
		InternalPort: "5432",
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
	externalPort, err := c.MappedPort(context.Background(), nat.Port(fmt.Sprintf("%s/tcp", pg.InternalPort)))
	if err != nil {
		return err
	}
	pg.ExternalPort = externalPort.Port()

	internalUrl, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User, pg.Password, pg.ContainerName, "5432", pg.DbName))
	if err != nil {
		return errors.Wrapf(err, "error parsing mercury db internal url")
	}
	pg.InternalURL = internalUrl
	externalUrl, err := url.Parse(fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		pg.User, pg.Password, externalPort.Port(), pg.DbName))
	if err != nil {
		return errors.Wrapf(err, "error parsing mercury db external url")
	}
	pg.ExternalURL = externalUrl

	pg.l.Info().
		Str("containerName", pg.ContainerName).
		Str("internalPort", pg.InternalPort).
		Str("externalPort", pg.ExternalPort).
		Str("internalURL", pg.InternalURL.String()).
		Str("externalURL", pg.ExternalURL.String()).
		Msg("Started Postgres DB container")

	return nil
}

func (pg *PostgresDb) ExecPgDump(stdout io.Writer) error {
	cmd := exec.Command("pg_dump", "-U", pg.User, "-h", "localhost", "-p", pg.ExternalPort, pg.DbName) //nolint:gosec
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", pg.Password),
	}
	cmd.Stdout = stdout

	return cmd.Run()
}

func (pg *PostgresDb) getContainerRequest() *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         pg.ContainerName,
		Image:        fmt.Sprintf("postgres:%s", pg.ImageVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", pg.InternalPort)},
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
