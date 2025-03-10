package test_env

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const defaultZookeeperImage = "confluentinc/cp-zookeeper:7.4.0"

type Zookeeper struct {
	EnvComponent
	InternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

// NewZookeeper creates a new Zookeeper instance with a unique container name and specified networks.
// It initializes default hooks and sets a startup timeout, making it suitable for managing distributed systems.
func NewZookeeper(networks []string) *Zookeeper {
	id, _ := uuid.NewRandom()
	z := &Zookeeper{
		EnvComponent: EnvComponent{
			ContainerName:  fmt.Sprintf("zookeper-%s", id.String()),
			Networks:       networks,
			StartupTimeout: 1 * time.Minute,
		},
	}
	z.SetDefaultHooks()
	return z
}

// WithTestInstance configures the Zookeeper instance for testing by setting up a test logger.
// It returns the modified Zookeeper instance, allowing for easier testing and logging during test execution.
func (z *Zookeeper) WithTestInstance(t *testing.T) *Zookeeper {
	z.l = logging.GetTestLogger(t)
	z.t = t
	return z
}

// WithContainerName sets the container name for the Zookeeper instance.
// It returns the updated Zookeeper instance, allowing for method chaining.
func (z *Zookeeper) WithContainerName(name string) *Zookeeper {
	z.ContainerName = name
	return z
}

// StartContainer initializes and starts a Zookeeper container.
// It configures logging, handles errors, and sets the internal URL for the container.
// This function is essential for setting up a Zookeeper instance in a test environment.
func (z *Zookeeper) StartContainer() error {
	l := logging.GetTestContainersGoTestLogger(z.t)
	cr, err := z.getContainerRequest()
	if err != nil {
		return err
	}
	req := tc.GenericContainerRequest{
		ContainerRequest: cr,
		Started:          true,
		Reuse:            true,
		Logger:           l,
	}
	c, err := tc.GenericContainer(testcontext.Get(z.t), req)
	if err != nil {
		return fmt.Errorf("cannot start Zookeeper container: %w", err)
	}
	name, err := c.Name(testcontext.Get(z.t))
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", -1)
	z.InternalUrl = fmt.Sprintf("%s:%s", name, "2181")

	z.l.Info().Str("containerName", name).
		Str("internalUrl", z.InternalUrl).
		Msgf("Started Zookeeper container")

	z.Container = c

	return nil
}

func (z *Zookeeper) getContainerRequest() (tc.ContainerRequest, error) {
	zookeeperImage := mirror.AddMirrorToImageIfSet(defaultZookeeperImage)
	return tc.ContainerRequest{
		Name:         z.ContainerName,
		Image:        zookeeperImage,
		ExposedPorts: []string{"2181/tcp"},
		Env: map[string]string{
			"ZOOKEEPER_CLIENT_PORT": "2181",
			"ZOOKEEPER_TICK_TIME":   "2000",
		},
		Networks: z.Networks,
		WaitingFor: tcwait.ForLog("ZooKeeper audit is disabled.").
			WithStartupTimeout(z.StartupTimeout).
			WithPollInterval(100 * time.Millisecond),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: z.PostStartsHooks,
				PostStops:  z.PostStopsHooks,
			},
		},
	}, nil
}
