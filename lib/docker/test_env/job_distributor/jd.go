package job_distributor

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const (
	JDContainerName            string = "job-distributor"
	DEAFULTJDContainerPort     string = "42242"
	DEFAULTCSAKeyEncryptionKey string = "!PASsword000!"
	DEAFULTWSRPCContainerPort  string = "8080"
)

type Option = func(j *Component)

type Component struct {
	test_env.EnvComponent
	Grpc                string
	Wsrpc               string
	InternalGRPC        string
	InternalWSRPC       string
	l                   zerolog.Logger
	t                   *testing.T
	dbConnection        string
	containerPort       string
	wsrpcPort           string
	csaKeyEncryptionKey string
}

func (j *Component) startOrRestartContainer(withReuse bool) error {
	req := j.getContainerRequest()
	l := logging.GetTestContainersGoTestLogger(j.t)
	c, err := docker.StartContainerWithRetry(j.l, tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
		Reuse:            withReuse,
		Logger:           l,
	})
	if err != nil {
		return err
	}
	j.Container = c
	ctx := testcontext.Get(j.t)
	host, err := test_env.GetHost(ctx, c)
	if err != nil {
		return errors.Wrapf(err, "cannot get host for container %s", j.ContainerName)
	}

	p, err := c.MappedPort(ctx, test_env.NatPort(j.containerPort))
	if err != nil {
		return errors.Wrapf(err, "cannot get container mapped port for container %s", j.ContainerName)
	}
	j.Grpc = fmt.Sprintf("%s:%s", host, p.Port())

	p, err = c.MappedPort(ctx, test_env.NatPort(j.wsrpcPort))
	if err != nil {
		return errors.Wrapf(err, "cannot get wsrpc mapped port for container %s", j.ContainerName)
	}
	j.Wsrpc = fmt.Sprintf("%s:%s", host, p.Port())
	j.InternalGRPC = fmt.Sprintf("%s:%s", j.ContainerName, j.containerPort)

	j.InternalWSRPC = fmt.Sprintf("%s:%s", j.ContainerName, j.wsrpcPort)
	j.l.Info().
		Str("containerName", j.ContainerName).
		Str("grpcURI", j.Grpc).
		Str("wsrpcURI", j.Wsrpc).
		Str("internalGRPC", j.InternalGRPC).
		Str("internalWSRPC", j.InternalWSRPC).
		Msg("Started Job Distributor container")

	return nil
}

func (j *Component) getContainerRequest() *tc.ContainerRequest {
	return &tc.ContainerRequest{
		Name:  j.ContainerName,
		Image: fmt.Sprintf("%s:%s", j.ContainerImage, j.ContainerVersion),
		ExposedPorts: []string{
			test_env.NatPortFormat(j.containerPort),
			test_env.NatPortFormat(j.wsrpcPort),
		},
		Env: map[string]string{
			"DATABASE_URL":              j.dbConnection,
			"PORT":                      j.containerPort,
			"NODE_RPC_PORT":             j.wsrpcPort,
			"CSA_KEY_ENCRYPTION_SECRET": j.csaKeyEncryptionKey,
		},
		Networks: j.Networks,
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(test_env.NatPort(j.containerPort)),
			tcwait.ForListeningPort(test_env.NatPort(j.wsrpcPort)),
		),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: j.PostStartsHooks,
				PostStops:  j.PostStopsHooks,
			},
		},
	}
}

func (j *Component) StartContainer() error {
	return j.startOrRestartContainer(false)
}

func (j *Component) RestartContainer() error {
	return j.startOrRestartContainer(true)
}

func New(networks []string, opts ...Option) *Component {
	id, _ := uuid.NewRandom()
	j := &Component{
		EnvComponent: test_env.EnvComponent{
			ContainerName:  fmt.Sprintf("%s-%s", JDContainerName, id.String()[0:8]),
			Networks:       networks,
			StartupTimeout: 2 * time.Minute,
		},
		containerPort:       DEAFULTJDContainerPort,
		wsrpcPort:           DEAFULTWSRPCContainerPort,
		csaKeyEncryptionKey: DEFAULTCSAKeyEncryptionKey,
		l:                   log.Logger,
	}
	j.SetDefaultHooks()
	for _, opt := range opts {
		opt(j)
	}
	return j
}

func WithTestInstance(t *testing.T) Option {
	return func(j *Component) {
		j.l = logging.GetTestLogger(t)
		j.t = t
	}
}

func WithContainerPort(port string) Option {
	return func(j *Component) {
		j.containerPort = port
	}
}

func WithWSRPCContainerPort(port string) Option {
	return func(j *Component) {
		j.wsrpcPort = port
	}
}

func WithDBURL(db string) Option {
	return func(j *Component) {
		if db != "" {
			j.dbConnection = db
		}
	}
}

func WithContainerName(name string) Option {
	return func(j *Component) {
		j.ContainerName = name
	}
}

func WithImage(image string) Option {
	return func(j *Component) {
		j.ContainerImage = image
	}
}

func WithVersion(version string) Option {
	return func(j *Component) {
		j.ContainerVersion = version
	}
}

func WithCSAKeyEncryptionKey(key string) Option {
	return func(j *Component) {
		j.csaKeyEncryptionKey = key
	}
}
