package test_env

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/imdario/mergo"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

type Kafka struct {
	EnvComponent
	Topics             []string
	BootstrapServerUrl string
	InternalUrl        string
	ExternalUrl        string
	l                  zerolog.Logger
	t                  *testing.T
}

func NewKafka(networks, topics []string) *Kafka {
	id, _ := uuid.NewRandom()
	return &Kafka{
		Topics: topics,
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("kafka-%s", id.String()),
			Networks:      networks,
		},
		l: log.Logger,
	}
}

func (k *Kafka) WithTestLogger(t *testing.T) *Kafka {
	k.l = logging.GetTestLogger(t)
	k.t = t
	return k
}

func (k *Kafka) WithContainerName(name string) *Kafka {
	k.ContainerName = name
	return k
}

func (k *Kafka) StartContainer(envVars map[string]string) error {
	k.InternalUrl = fmt.Sprintf("%s:%s", k.ContainerName, "9092")
	// TODO: Fix mapped port
	k.ExternalUrl = fmt.Sprintf("localhost:%s", "29092")
	k.BootstrapServerUrl = k.InternalUrl

	l := tc.Logger
	if k.t != nil {
		l = logging.CustomT{
			T: k.t,
			L: k.l,
		}
	}
	req := tc.GenericContainerRequest{
		ContainerRequest: k.getContainerRequest(envVars),
		Started:          true,
		Reuse:            true,
		Logger:           l,
	}
	c, err := tc.GenericContainer(context.Background(), req)
	if err != nil {
		return errors.Wrapf(err, "cannot start Kafka container")
	}

	k.l.Info().Str("containerName", k.ContainerName).
		Str("internalUrl", k.InternalUrl).
		Str("externalUrl", k.ExternalUrl).
		Msgf("Started Kafka container")

	k.Container = c

	return nil
}

func (k *Kafka) CreateLocalTopics() error {
	for _, topic := range k.Topics {
		cmd := []string{"kafka-topics", "--bootstrap-server", fmt.Sprintf("http://%s", k.BootstrapServerUrl),
			"--topic", topic, "--create", "--if-not-exists", "--partitions", "25",
			"--replication-factor", "1"}
		code, output, err := k.Container.Exec(context.Background(), cmd)
		if err != nil {
			return err
		}
		if code != 0 {
			outputBytes, _ := io.ReadAll(output)
			outputString := strings.TrimSpace(string(outputBytes))
			return errors.Errorf("Create topics returned %d code. Output: %s", code, outputString)
		}
		k.l.Info().
			Strs("cmd", cmd).
			Msgf("Created Kafka %s topic", topic)
	}
	return nil
}

func (k *Kafka) getContainerRequest(envVars map[string]string) tc.ContainerRequest {
	defaultValues := map[string]string{
		"KAFKA_BROKER_ID":                                "1",
		"KAFKA_ZOOKEEPER_CONNECT":                        "zookeeper:2181",
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
		"KAFKA_ADVERTISED_LISTENERS":                     fmt.Sprintf("PLAINTEXT://%s,PLAINTEXT_HOST://%s", k.InternalUrl, k.ExternalUrl),
		"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
		"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS":         "0",
		"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
		"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
		"KAFKA_CREATE_TOPICS":                            "reports_instant:1:1,reports_dlq:1:1",
	}
	if err := mergo.Merge(defaultValues, envVars, mergo.WithOverride); err != nil {
		panic(err)
	}
	return tc.ContainerRequest{
		Name:         k.ContainerName,
		Image:        "confluentinc/cp-kafka:7.4.0",
		ExposedPorts: []string{"29092/tcp"},
		Env:          defaultValues,
		Networks:     k.Networks,
		WaitingFor: tcwait.ForLog("[KafkaServer id=1] started (kafka.server.KafkaServer)").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}
}
