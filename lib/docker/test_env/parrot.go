package test_env

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const defaultParrotImage = "parrot:latest"

// Parrot is a test environment component that wraps a Parrot server.
type Parrot struct {
	EnvComponent
	Client           *parrot.Client
	ExternalEndpoint string
	InternalEndpoint string
	t                *testing.T
	l                zerolog.Logger
}

// NewParrot creates a new instance of ParrotServer with specified networks and options.
// It initializes the server with a unique container name and a default startup timeout.
// This function is useful for testing decentralized applications in a controlled environment.
func NewParrot(networks []string, opts ...EnvComponentOption) *Parrot {
	p := &Parrot{
		EnvComponent: EnvComponent{
			ContainerName:  fmt.Sprintf("%s-%s", "parrot", uuid.NewString()[0:8]),
			Networks:       networks,
			StartupTimeout: 1 * time.Minute,
		},
		l: log.Logger,
	}
	for _, opt := range opts {
		opt(&p.EnvComponent)
	}
	return p
}

// WithTestInstance configures the MockServer with a test logger and test context.
// It returns the updated MockServer instance for use in testing scenarios.
func (p *Parrot) WithTestInstance(t *testing.T) *Parrot {
	p.l = logging.GetTestLogger(t)
	p.t = t
	return p
}

// SetExternalAdapterMocks configures a specified number of mock external adapter endpoints.
// It generates unique paths for each adapter and stores their URLs for later use.
// This function is useful for testing scenarios that require multiple external adapter interactions.
func (p *Parrot) SetExternalAdapterMocks(count int) error {
	// for i := 0; i < count; i++ {
	// 	path := fmt.Sprintf("/ea-%d", i)
	// 	err := ms.Client.SetRandomValuePath(path)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	cName, err := ms.Container.Name(testcontext.Get(ms.t))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	cName = strings.Replace(cName, "/", "", -1)
	// 	eaUrl, err := url.Parse(fmt.Sprintf("http://%s:%s%s",
	// 		cName, "1080", path))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	ms.EAMockUrls = append(ms.EAMockUrls, eaUrl)
	// }
	// return nil
	return fmt.Errorf("not implemented")
}

// StartContainer initializes and starts a Parrot container.
// It sets up logging, retrieves the container request, and establishes endpoints for communication.
// This function is essential for testing environments that require a mock server instance.
func (p *Parrot) StartContainer() error {
	l := logging.GetTestContainersGoTestLogger(p.t)
	cr, err := p.getContainerRequest()
	if err != nil {
		return err
	}
	c, err := docker.StartContainerWithRetry(p.l, tc.GenericContainerRequest{
		ContainerRequest: cr,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start MockServer container: %w", err)
	}
	p.Container = c
	endpoint, err := GetEndpoint(testcontext.Get(p.t), c, "http")
	if err != nil {
		return err
	}
	p.ExternalEndpoint = endpoint
	p.InternalEndpoint = fmt.Sprintf("http://%s", p.ContainerName)

	p.Client = parrot.NewClient(p.ExternalEndpoint)

	p.l.Info().Str("External Endpoint", p.ExternalEndpoint).
		Str("Internal Endpoint", p.InternalEndpoint).
		Str("Container Name", p.ContainerName).
		Msg("Started Parrot Container")
	return nil
}

func (p *Parrot) getContainerRequest() (tc.ContainerRequest, error) {
	pImage := mirror.AddMirrorToImageIfSet(defaultParrotImage)

	return tc.ContainerRequest{
		Name:         p.ContainerName,
		Image:        pImage,
		ExposedPorts: []string{"80/tcp"},
		Networks:     p.Networks,
		WaitingFor: tcwait.ForHealthCheck().
			WithPollInterval(100 * time.Millisecond).WithStartupTimeout(p.StartupTimeout),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: p.PostStartsHooks,
				PostStops:  p.PostStopsHooks,
			},
		},
	}, nil
}
