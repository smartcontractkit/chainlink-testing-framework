package chipingressset

import (
	"context"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func CreateTopics(ctx context.Context, kafkaAddress string, topics []string) error {
	if len(topics) == 0 {
		return nil
	}
	framework.L.Debug().Msgf("Creating Kafka topics: %v", topics)

	admin, adminErr := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAddress,
	})
	if adminErr != nil {
		return errors.Wrap(adminErr, "failed to create kafka admin client")
	}
	defer admin.Close()

	for _, topic := range topics {
		spec := kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		_, resultsErr := admin.CreateTopics(ctx, []kafka.TopicSpecification{spec}, kafka.SetAdminOperationTimeout(5000))
		if resultsErr != nil {
			return errors.Wrapf(resultsErr, "failed to create topic %s", topic)
		}
	}
	framework.L.Debug().Msgf("Created Kafka %d topics", len(topics))

	return nil
}

func DeleteAllTopics(ctx context.Context, kafkaAddress string) error {
	framework.L.Debug().Msg("Deleting all Kafka topics")

	admin, adminErr := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAddress,
	})
	if adminErr != nil {
		return errors.Wrap(adminErr, "failed to create kafka admin client")
	}
	defer admin.Close()

	// Get metadata for all topics
	metadata, metaErr := admin.GetMetadata(nil, false, 5000)
	if metaErr != nil {
		return errors.Wrap(metaErr, "failed to fetch metadata")
	}

	// Collect all topic names, skipping internal topics (e.g., __consumer_offsets)
	var topicNames []string
	for topicName := range metadata.Topics {
		if !strings.HasPrefix(topicName, "__") {
			topicNames = append(topicNames, topicName)
		}
	}

	if len(topicNames) == 0 {
		framework.L.Debug().Msg("No topics to delete")
		return nil
	}

	// Delete all collected topics
	_, deleteErr := admin.DeleteTopics(ctx, topicNames, kafka.SetAdminOperationTimeout(5000))
	if deleteErr != nil {
		return errors.Wrap(deleteErr, "failed to delete topics")
	}

	framework.L.Debug().Msgf("Deleted %d topics", len(topicNames))

	return nil
}
