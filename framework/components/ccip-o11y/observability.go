package fake

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgreSQLInput struct {
	Image string `toml:"image"`
	Port  string `toml:"port"`
}

type PostgreSQLOutput struct {
	URL string `toml:"url"`
}

type KafkaOutput struct {
	URL string `toml:"url"`
}

type KafkaInput struct {
	Image          string `toml:"image"`
	PlaintextPort  string `toml:"plaintext_port"`
	BrokerPort     string `toml:"broker_port"`
	ControllerPort string `toml:"controller_port"`
}

type Input struct {
	PostgreSQL *PostgreSQLInput `toml:"postgresql"`
	Kafka      *KafkaInput      `toml:"kafka"`
	Out        *Output          `toml:"out"`
}

type Output struct {
	UseCache   bool              `toml:"use_cache"`
	PostgreSQL *PostgreSQLOutput `toml:"postgresql"`
	Kafka      *KafkaOutput      `toml:"kafka"`
}

// NewCCIPO11y provides services for CCIP o11y - PostgreSQL and Kafka
func NewCCIPO11y(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	ctx := context.Background()

	_, err := postgres.Run(ctx,
		in.PostgreSQL.Image,
		tc.WithHostConfigModifier(func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", in.PostgreSQL.Port)): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%s/tcp", in.PostgreSQL.Port),
					},
				},
			}
		}),
		tc.WithLabels(map[string]string{"framework": "ctf"}),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run postgresql: %w", err)
	}

	ptPort := fmt.Sprintf("%s/tcp", in.Kafka.PlaintextPort)
	brokerPort := fmt.Sprintf("%s/tcp", in.Kafka.BrokerPort)
	controllerPort := fmt.Sprintf("%s/tcp", in.Kafka.ControllerPort)

	_, err = kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		tc.WithExposedPorts(ptPort, brokerPort, controllerPort),
		tc.WithHostConfigModifier(func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				nat.Port(ptPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: ptPort,
					},
				},
				nat.Port(brokerPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: brokerPort,
					},
				},
				nat.Port(controllerPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: controllerPort,
					},
				},
			}
		}),
		tc.WithLabels(map[string]string{"framework": "ctf"}),
		kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run kafka: %w", err)
	}
	out := &Output{
		PostgreSQL: &PostgreSQLOutput{
			URL: fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", "test", "test", in.PostgreSQL.Port, "test"),
		},
		Kafka: &KafkaOutput{
			URL: fmt.Sprintf("http://localhost:%s", in.Kafka.PlaintextPort),
		},
	}
	in.Out = out
	return out, nil
}
