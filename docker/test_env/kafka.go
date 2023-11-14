package test_env

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/imdario/mergo"

	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

type Kafka struct {
	EnvComponent
	TopicConfigs       []KafkaTopicConfig
	EnvVars            map[string]string
	BootstrapServerUrl string
	InternalUrl        string
	ExternalUrl        string
	l                  zerolog.Logger
	t                  *testing.T
}

type KafkaTopicConfig struct {
	TopicName     string `json:"topic_name"`
	Partitions    int    `json:"partitions"`
	Replication   int    `json:"replication"`
	CleanupPolicy string `json:"cleanup_policy"`
}

func NewKafka(networks []string) *Kafka {
	id, _ := uuid.NewRandom()
	containerName := fmt.Sprintf("kafka-%s", id.String())
	defaultEnvVars := map[string]string{
		"KAFKA_BROKER_ID":                                "1",
		"KAFKA_ZOOKEEPER_CONNECT":                        "zookeeper:2181",
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
		"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
		"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS":         "0",
		"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
		"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
		"KAFKA_CREATE_TOPICS":                            "reports_instant:1:1,reports_dlq:1:1",
	}
	return &Kafka{
		EnvComponent: EnvComponent{
			ContainerName: containerName,
			Networks:      networks,
		},
		EnvVars: defaultEnvVars,
		l:       log.Logger,
	}
}

func (k *Kafka) WithTestLogger(t *testing.T) *Kafka {
	k.l = logging.GetTestLogger(t)
	k.t = t
	return k
}

func (k *Kafka) WithContainerName(name string) *Kafka {
	k.ContainerName = name
	internalUrl := fmt.Sprintf("%s:%s", name, "9092")
	bootstrapServerUrl := internalUrl
	k.InternalUrl = internalUrl
	k.BootstrapServerUrl = bootstrapServerUrl
	return k
}

func (k *Kafka) WithTopics(topics []KafkaTopicConfig) *Kafka {
	k.TopicConfigs = topics
	return k
}

func (k *Kafka) WithZookeeper(zookeeperUrl string) *Kafka {
	envVars := map[string]string{
		"KAFKA_ZOOKEEPER_CONNECT": zookeeperUrl,
	}
	return k.WithEnvVars(envVars)
}

func (k *Kafka) WithEnvVars(envVars map[string]string) *Kafka {
	if err := mergo.Merge(&k.EnvVars, envVars, mergo.WithOverride); err != nil {
		k.l.Fatal().Err(err).Msg("Failed to merge env vars")
	}
	return k
}

func (k *Kafka) StartContainer() error {
	l := logging.GetTestContainersGoTestLogger(k.t)
	k.InternalUrl = fmt.Sprintf("%s:%s", k.ContainerName, "9092")
	// TODO: Fix mapped port
	k.ExternalUrl = fmt.Sprintf("127.0.0.1:%s", "29092")
	k.BootstrapServerUrl = k.InternalUrl
	envVars := map[string]string{
		"KAFKA_ADVERTISED_LISTENERS": fmt.Sprintf("PLAINTEXT://%s,PLAINTEXT_HOST://%s", k.InternalUrl, k.ExternalUrl),
	}
	k.WithEnvVars(envVars)
	cr, err := k.getContainerRequest()
	if err != nil {
		return err
	}
	c, err := docker.StartContainerWithRetry(k.l, tc.GenericContainerRequest{
		ContainerRequest: cr,
		Started:          true,
		Reuse:            true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start Kafka container: %w", err)
	}

	k.l.Info().Str("containerName", k.ContainerName).
		Str("internalUrl", k.InternalUrl).
		Str("externalUrl", k.ExternalUrl).
		Msgf("Started Kafka container")

	k.Container = c

	return nil
}

func (k *Kafka) CreateLocalTopics() error {
	for _, topicConfig := range k.TopicConfigs {
		cmd := []string{"kafka-topics", "--bootstrap-server", fmt.Sprintf("http://%s", k.BootstrapServerUrl),
			"--topic", topicConfig.TopicName,
			"--create",
			"--if-not-exists",
			"--partitions", fmt.Sprintf("%d", topicConfig.Partitions),
			"--replication-factor", fmt.Sprintf("%d", topicConfig.Replication)}
		if topicConfig.CleanupPolicy != "" {
			cmd = append(cmd, "--config", fmt.Sprintf("cleanup.policy=%s", topicConfig.CleanupPolicy))
		}
		code, output, err := k.Container.Exec(utils.TestContext(k.t), cmd)
		if err != nil {
			return err
		}
		if code != 0 {
			outputBytes, _ := io.ReadAll(output)
			outputString := strings.TrimSpace(string(outputBytes))
			return fmt.Errorf("Create topics returned %d code. Output: %s", code, outputString)
		}
		k.l.Info().
			Strs("cmd", cmd).
			Msgf("Created Kafka %s topic with partitions: %d, replication: %d, cleanup.policy: %s",
				topicConfig.TopicName, topicConfig.Partitions, topicConfig.Replication, topicConfig.CleanupPolicy)
	}
	return nil
}

func (k *Kafka) getContainerRequest() (tc.ContainerRequest, error) {
	kafkaImage, err := mirror.GetImage("confluentinc/cp-kafka")
	if err != nil {
		return tc.ContainerRequest{}, err
	}
	return tc.ContainerRequest{
		Name:         k.ContainerName,
		Image:        kafkaImage,
		ExposedPorts: []string{"29092/tcp"},
		Env:          k.EnvVars,
		Networks:     k.Networks,
		WaitingFor: tcwait.ForLog("[KafkaServer id=1] started (kafka.server.KafkaServer)").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}, nil
}
