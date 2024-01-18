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

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

type Zookeeper struct {
	EnvComponent
	InternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

func NewZookeeper(networks []string) *Zookeeper {
	id, _ := uuid.NewRandom()
	z := &Zookeeper{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("zookeper-%s", id.String()),
			Networks:      networks,
		},
	}
	z.SetDefaultHooks()
	return z
}

func (z *Zookeeper) WithTestInstance(t *testing.T) *Zookeeper {
	z.l = logging.GetTestLogger(t)
	z.t = t
	return z
}

func (z *Zookeeper) WithContainerName(name string) *Zookeeper {
	z.ContainerName = name
	return z
}

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
		return fmt.Errorf("cannot start Zookeper container: %w", err)
	}
	name, err := c.Name(testcontext.Get(z.t))
	if err != nil {
		return err
	}
	name = strings.Replace(name, "/", "", -1)
	z.InternalUrl = fmt.Sprintf("%s:%s", name, "2181")

	z.l.Info().Str("containerName", name).
		Str("internalUrl", z.InternalUrl).
		Msgf("Started Zookeper container")

	z.Container = c

	return nil
}

func (z *Zookeeper) getContainerRequest() (tc.ContainerRequest, error) {
	zookeeperImage, err := mirror.GetImage("confluentinc/cp-zookeeper")
	if err != nil {
		return tc.ContainerRequest{}, err
	}
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
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
		LifecycleHooks: []tc.ContainerLifecycleHooks{
			{
				PostStarts: z.PostStartsHooks,
				PostStops:  z.PostStopsHooks,
			},
		},
	}, nil
}
