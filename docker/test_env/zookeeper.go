package test_env

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type Zookeeper struct {
	EnvComponent
	InternalUrl string
	l           zerolog.Logger
	t           *testing.T
}

func NewZookeeper(networks []string) *Zookeeper {
	id, _ := uuid.NewRandom()
	return &Zookeeper{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("zookeper-%s", id.String()),
			Networks:      networks,
		},
	}
}

func (z *Zookeeper) WithTestLogger(t *testing.T) *Zookeeper {
	z.l = logging.GetTestLogger(t)
	z.t = t
	return z
}

func (z *Zookeeper) WithContainerName(name string) *Zookeeper {
	z.ContainerName = name
	return z
}

func (z *Zookeeper) StartContainer() error {
	l := tc.Logger
	if z.t != nil {
		l = logging.CustomT{
			T: z.t,
			L: z.l,
		}
	}
	req := tc.GenericContainerRequest{
		ContainerRequest: z.getContainerRequest(),
		Started:          true,
		Reuse:            true,
		Logger:           l,
	}
	c, err := tc.GenericContainer(context.Background(), req)
	if err != nil {
		return errors.Wrapf(err, "cannot start Zookeper container")
	}
	name, err := c.Name(context.Background())
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

func (z *Zookeeper) getContainerRequest() tc.ContainerRequest {
	return tc.ContainerRequest{
		Name:         z.ContainerName,
		Image:        "confluentinc/cp-zookeeper:7.4.0",
		ExposedPorts: []string{"2181/tcp"},
		Env: map[string]string{
			"ZOOKEEPER_CLIENT_PORT": "2181",
			"ZOOKEEPER_TICK_TIME":   "2000",
		},
		Networks: z.Networks,
		WaitingFor: tcwait.ForLog("ZooKeeper audit is disabled.").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
