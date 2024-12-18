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

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const defaultKafkaImage = "confluentinc/cp-kafka:7.4.0"

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

// NewKafka initializes a new Kafka instance with a unique container name and default environment variables.
// It sets up the necessary configurations for Kafka to operate within specified networks, making it easy to deploy and manage Kafka services in a containerized environment.
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
			ContainerName:  containerName,
			Networks:       networks,
			StartupTimeout: 1 * time.Minute,
		},
		EnvVars: defaultEnvVars,
		l:       log.Logger,
	}
}

// WithTestInstance configures the Kafka instance for testing by setting up a test logger.
// It returns the modified Kafka instance, allowing for easier testing and logging during unit tests.
func (k *Kafka) WithTestInstance(t *testing.T) *Kafka {
	k.l = logging.GetTestLogger(t)
	k.t = t
	return k
}

// WithContainerName sets the container name for the Kafka instance.
// It configures the internal and bootstrap server URLs based on the provided name,
// allowing for easy identification and connection to the Kafka service.
func (k *Kafka) WithContainerName(name string) *Kafka {
	k.ContainerName = name
	internalUrl := fmt.Sprintf("%s:%s", name, "9092")
	bootstrapServerUrl := internalUrl
	k.InternalUrl = internalUrl
	k.BootstrapServerUrl = bootstrapServerUrl
	return k
}

// WithTopics sets the Kafka topic configurations for the Kafka instance.
// It returns the updated Kafka instance, allowing for method chaining.
func (k *Kafka) WithTopics(topics []KafkaTopicConfig) *Kafka {
	k.TopicConfigs = topics
	return k
}

// WithZookeeper sets the Zookeeper connection URL for the Kafka instance.
// It prepares the necessary environment variables and returns the updated Kafka instance.
func (k *Kafka) WithZookeeper(zookeeperUrl string) *Kafka {
	envVars := map[string]string{
		"KAFKA_ZOOKEEPER_CONNECT": zookeeperUrl,
	}
	return k.WithEnvVars(envVars)
}

// WithEnvVars merges the provided environment variables into the Kafka instance's existing environment variables.
// It allows customization of the Kafka container's configuration before starting it.
func (k *Kafka) WithEnvVars(envVars map[string]string) *Kafka {
	if err := mergo.Merge(&k.EnvVars, envVars, mergo.WithOverride); err != nil {
		k.l.Fatal().Err(err).Msg("Failed to merge env vars")
	}
	return k
}

// StartContainer initializes and starts a Kafka container with specified environment variables.
// It sets internal and external URLs for the container and logs the startup process.
// This function is essential for setting up a Kafka instance for testing or development purposes.
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

// CreateLocalTopics creates Kafka topics based on the provided configurations.
// It ensures that topics are only created if they do not already exist,
// and logs the creation details for each topic. This function is useful
// for initializing Kafka environments with predefined topic settings.
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
		code, output, err := k.Container.Exec(testcontext.Get(k.t), cmd)
		if err != nil {
			return err
		}
		if code != 0 {
			outputBytes, _ := io.ReadAll(output)
			outputString := strings.TrimSpace(string(outputBytes))
			return fmt.Errorf("create topics returned %d code. Output: %s", code, outputString)
		}
		k.l.Info().
			Strs("cmd", cmd).
			Msgf("Created Kafka %s topic with partitions: %d, replication: %d, cleanup.policy: %s",
				topicConfig.TopicName, topicConfig.Partitions, topicConfig.Replication, topicConfig.CleanupPolicy)
	}
	return nil
}

func (k *Kafka) getContainerRequest() (tc.ContainerRequest, error) {
	kafkaImage := mirror.AddMirrorToImageIfSet(defaultKafkaImage)
	return tc.ContainerRequest{
		Name:         k.ContainerName,
		Image:        kafkaImage,
		ExposedPorts: []string{"29092/tcp"},
		Env:          k.EnvVars,
		Networks:     k.Networks,
		WaitingFor: tcwait.ForLog("[KafkaServer id=1] started (kafka.server.KafkaServer)").
			WithStartupTimeout(k.StartupTimeout).
			WithPollInterval(100 * time.Millisecond),
	}, nil
}
