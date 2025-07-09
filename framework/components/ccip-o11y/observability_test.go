package fake

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCCIPo11y(t *testing.T) {
	out, err := NewCCIPO11y(&Input{
		PostgreSQL: &PostgreSQLInput{
			Image: "postgres:16-alpine",
			Port:  "5432",
		},
		Kafka: &KafkaInput{
			Image:          "confluentinc/confluent-local:7.5.0",
			PlaintextPort:  "9093",
			BrokerPort:     "9092",
			ControllerPort: "9094",
		},
	})
	require.NoError(t, err)
	// psql "postgres://test:test@localhost:5432/test?sslmode=disable"
	fmt.Printf("PostgreSQL URL: %s\n", out.PostgreSQL.URL)
	// brew install kafka
	// kafka-topics --create --bootstrap-server localhost:9093 --replication-factor 1 --partitions 1 --topic my-test-topic
	// kafka-console-producer --bootstrap-server localhost:9093 --topic my-test-topic
	// kafka-console-consumer --bootstrap-server localhost:9093 --topic my-test-topic
	fmt.Printf("Kafka URL: %s\n", out.Kafka.URL)
}
