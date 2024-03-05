package test_env

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

type TestLogConsumer struct {
	Msgs []string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.Msgs = append(g.Msgs, string(l.Content))
}

func followLogs(t *testing.T, c testcontainers.Container) *TestLogConsumer {
	consumer := &TestLogConsumer{
		Msgs: make([]string, 0),
	}
	go func() {
		c.FollowOutput(consumer)
		err := c.StartLogProducer(testcontext.Get(t), testcontainers.WithLogProductionTimeout(time.Duration(5*time.Second)))
		require.NoError(t, err)
	}()
	return consumer
}

func TestEnvComponentPauseChaos(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)

	defaultChainCfg := GetDefaultChainConfig()
	g := NewGethPoa([]string{network.Name}, &defaultChainCfg).
		WithTestInstance(t)
	_, err = g.StartContainer()
	require.NoError(t, err)
	t.Run("check that testcontainers can be paused", func(t *testing.T) {
		consumer := followLogs(t, *g.GetContainer())

		timeStrNow := time.Now().Add(3 * time.Second).UTC().String()
		justTime := strings.Split(timeStrNow, " ")[1]
		justTimeWithoutMicrosecs := justTime[:len(justTime)-7]

		// blocking
		err = g.(*Geth).ChaosPause(l, 5*time.Second)

		// check that there were no logs when paused
		for _, lo := range consumer.Msgs {
			if strings.Contains(lo, justTimeWithoutMicrosecs) {
				t.Fail()
			}
		}
	})

	t.Run("check container traffic can be lost", func(t *testing.T) {
		// TODO: assert with a busybox container that the traffic is lost
		err = g.(*Geth).ChaosNetworkLoss(l, 30*time.Second, 100, "", nil, nil, nil)
		require.NoError(t, err)
	})
	t.Run("check container latency can be changed", func(t *testing.T) {
		// TODO: assert with a busybox container that the traffic is delayed
		err = g.(*Geth).ChaosNetworkDelay(l, 30*time.Second, 5*time.Second, "", nil, nil, nil)
		require.NoError(t, err)
	})
}
