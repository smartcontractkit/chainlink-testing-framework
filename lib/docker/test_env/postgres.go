package test_env

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const defaultPostgresImage = "postgres:15.6"

type PostgresDb struct {
	EnvComponent
	User         string   `json:"user"`
	Password     string   `json:"password"`
	DbName       string   `json:"dbName"`
	InternalPort string   `json:"internalPort"`
	ExternalPort string   `json:"-"`
	InternalURL  *url.URL `json:"-"`
	ExternalURL  *url.URL `json:"-"`
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

func WithPostgresImageName(imageName string) PostgresDbOption {
	return func(c *PostgresDb) {
		if imageName != "" {
			c.ContainerImage = imageName
		}
	}
}

func WithPostgresDbLogStream(ls *logstream.LogStream) PostgresDbOption {
	return func(c *PostgresDb) {
		c.LogStream = ls
	}
}

func WithPostgresImageVersion(version string) PostgresDbOption {
	return func(c *PostgresDb) {
		if version != "" {
			c.ContainerVersion = version
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

func WithContainerEnv(key, value string) PostgresDbOption {
	return func(c *PostgresDb) {
		c.ContainerEnvs[key] = value
	}
}

func NewPostgresDb(networks []string, opts ...PostgresDbOption) (*PostgresDb, error) {
	imageParts := strings.Split(defaultPostgresImage, ":")
	pg := &PostgresDb{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "postgres-db", uuid.NewString()[0:8]),
			ContainerImage:   imageParts[0],
			ContainerVersion: imageParts[1],
			ContainerEnvs:    map[string]string{},
			Networks:         networks,
			StartupTimeout:   2 * time.Minute,
		},
		User:         "postgres",
		Password:     "mysecretpassword",
		DbName:       "testdb",
		InternalPort: "5432",
		l:            log.Logger,
	}

	pg.SetDefaultHooks()
	for _, opt := range opts {
		opt(pg)
	}

	// if the internal docker repo is set then add it to the version
	pg.EnvComponent.ContainerImage = mirror.AddMirrorToImageIfSet(pg.EnvComponent.ContainerImage)

	// Set default container envs
	pg.ContainerEnvs["POSTGRES_USER"] = pg.User
	pg.ContainerEnvs["POSTGRES_DB"] = pg.DbName
	pg.ContainerEnvs["POSTGRES_PASSWORD"] = pg.Password

	return pg, nil
}

func (pg *PostgresDb) WithTestInstance(t *testing.T) *PostgresDb {
	pg.l = logging.GetTestLogger(t)
	pg.t = t
	return pg
}

func (pg *PostgresDb) StartContainer() error {
	return pg.startOrRestartContainer(false)
}

func (pg *PostgresDb) RestartContainer() error {
	return pg.startOrRestartContainer(true)
}

func (pg *PostgresDb) startOrRestartContainer(withReuse bool) error {
	req := pg.getContainerRequest()
	l := logging.GetTestContainersGoTestLogger(pg.t)
	c, err := docker.StartContainerWithRetry(pg.l, tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
		Reuse:            withReuse,
		Logger:           l,
	})
	if err != nil {
		return err
	}
	pg.Container = c
	externalPort, err := c.MappedPort(testcontext.Get(pg.t), nat.Port(fmt.Sprintf("%s/tcp", pg.InternalPort)))
	if err != nil {
		return err
	}
	pg.ExternalPort = externalPort.Port()

	internalUrl, err := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pg.User, pg.Password, pg.ContainerName, "5432", pg.DbName))
	if err != nil {
		return fmt.Errorf("error parsing db internal url: %w", err)
	}
	pg.InternalURL = internalUrl
	externalUrl, err := url.Parse(fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
		pg.User, pg.Password, externalPort.Port(), pg.DbName))
	if err != nil {
		return fmt.Errorf("error parsing db external url: %w", err)
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

// ExecPgDumpFromLocal executes pg_dump from local machine by connecting to external Postgres port. For it to work pg_dump
// needs to be installed on local machine.
func (pg *PostgresDb) ExecPgDumpFromLocal(stdout io.Writer) error {
	cmd := exec.Command("pg_dump", "-U", pg.User, "-h", "127.0.0.1", "-p", pg.ExternalPort, pg.DbName) //nolint:gosec
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", pg.Password),
	}
	cmd.Stdout = stdout

	return cmd.Run()
}

// ExecPgDumpFromContainer executed pg_dump from inside the container. It dumps it to temporary file inside the container
// and then writes to the writer.
func (pg *PostgresDb) ExecPgDumpFromContainer(writer io.Writer) error {
	tmpFile := "/tmp/db_dump.sql"
	command := []string{"pg_dump", "-U", pg.User, "-f", tmpFile, pg.DbName}
	env := []string{
		fmt.Sprintf("PGPASSWORD=%s", pg.Password),
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Minute)

	_, _, err := pg.Container.Exec(ctx, command, tcexec.WithEnv(env))
	if err != nil {
		cancelFn()
		return errors.Wrap(err, "Failed to execute pg_dump")
	}

	reader, err := pg.Container.CopyFileFromContainer(ctx, tmpFile)
	if err != nil {
		cancelFn()
		return errors.Wrapf(err, "Failed to open for reading %s temporary file with db dump", tmpFile)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		cancelFn()
		return errors.Wrapf(err, "Failed to send data from %s temporary file with db dump", tmpFile)
	}

	cancelFn()

	return nil
}

func (pg *PostgresDb) getContainerRequest() *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:         pg.ContainerName,
		Image:        fmt.Sprintf("%s:%s", pg.ContainerImage, pg.ContainerVersion),
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", pg.InternalPort)},
		Env:          pg.ContainerEnvs,
		Networks:     pg.Networks,
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "127.0.0.1",
			"-U", pg.User, "-c", "select", "1", "-d", pg.DbName}).
			WithStartupTimeout(pg.StartupTimeout),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: pg.PostStartsHooks,
				PostStops:  pg.PostStopsHooks,
			},
		},
	}
}
