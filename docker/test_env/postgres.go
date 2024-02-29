package test_env

import (
	"fmt"
	"io"
	"net/url"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

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
	image, err := mirror.GetImage("postgres")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(image, ":")
	pg := &PostgresDb{
		EnvComponent: EnvComponent{
			ContainerName:    fmt.Sprintf("%s-%s", "postgres-db", uuid.NewString()[0:8]),
			ContainerImage:   parts[0],
			ContainerVersion: parts[1],
			ContainerEnvs:    map[string]string{},
			Networks:         networks,
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
	req := pg.getContainerRequest()
	l := logging.GetTestContainersGoTestLogger(pg.t)
	c, err := docker.StartContainerWithRetry(pg.l, tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
		Reuse:            true,
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
		return fmt.Errorf("error parsing mercury db internal url: %w", err)
	}
	pg.InternalURL = internalUrl
	externalUrl, err := url.Parse(fmt.Sprintf("postgres://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
		pg.User, pg.Password, externalPort.Port(), pg.DbName))
	if err != nil {
		return fmt.Errorf("error parsing mercury db external url: %w", err)
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
	cmd := exec.Command("pg_dump", "-U", pg.User, "-h", "127.0.0.1", "-p", pg.ExternalPort, pg.DbName) //nolint:gosec
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", pg.Password),
	}
	cmd.Stdout = stdout

	return cmd.Run()
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
			WithStartupTimeout(10 * time.Second),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: pg.PostStartsHooks,
				PostStops:  pg.PostStopsHooks,
			},
		},
	}
}
