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

// StartContainer initializes and starts a container for the component.
// It ensures the container is ready for use, logging relevant connection details.
// This function is essential for setting up the environment required for the component's operations.
func (j *Component) StartContainer() error {
	return j.startOrRestartContainer(false)
}

// RestartContainer restarts the container associated with the Component.
// It ensures that the container is started with the option to reuse resources,
// making it useful for maintaining service continuity during updates or failures.
func (j *Component) RestartContainer() error {
	return j.startOrRestartContainer(true)
}

// New creates a new Component instance with a unique container name and specified networks.
// It allows for optional configuration through functional options and sets default values for various parameters.
// This function is useful for initializing components in a test environment.
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

// WithTestInstance sets up a test logger and test context for a Component.
// It is useful for initializing components in unit tests, ensuring that logs
// are captured and associated with the provided testing.T instance.
func WithTestInstance(t *testing.T) Option {
	return func(j *Component) {
		j.l = logging.GetTestLogger(t)
		j.t = t
	}
}

// WithContainerPort sets the container port for a Component.
// This option allows customization of the network configuration for the container,
// enabling proper communication with external services.
func WithContainerPort(port string) Option {
	return func(j *Component) {
		j.containerPort = port
	}
}

// WithWSRPCContainerPort sets the WebSocket RPC port for the Component.
// This option allows users to configure the port for WebSocket communication,
// enabling flexibility in service deployment and integration.
func WithWSRPCContainerPort(port string) Option {
	return func(j *Component) {
		j.wsrpcPort = port
	}
}

// WithDBURL returns an Option that sets the database connection URL for a Component.
// It allows users to configure the Component with a specific database connection string.
func WithDBURL(db string) Option {
	return func(j *Component) {
		if db != "" {
			j.dbConnection = db
		}
	}
}

// WithContainerName sets the container name for a Component.
// It returns an Option that can be used to configure the Component during initialization.
func WithContainerName(name string) Option {
	return func(j *Component) {
		j.ContainerName = name
	}
}

// WithImage sets the container image for a Component.
// This option allows users to specify which image to use, enabling customization of the component's runtime environment.
func WithImage(image string) Option {
	return func(j *Component) {
		j.ContainerImage = image
	}
}

// WithVersion sets the container version for a Component.
// It allows users to specify the desired version, ensuring the Component runs with the correct configuration.
func WithVersion(version string) Option {
	return func(j *Component) {
		j.ContainerVersion = version
	}
}

// WithCSAKeyEncryptionKey sets the CSA key encryption key for a Component.
// This function is useful for configuring secure key management in decentralized applications.
func WithCSAKeyEncryptionKey(key string) Option {
	return func(j *Component) {
		j.csaKeyEncryptionKey = key
	}
}
