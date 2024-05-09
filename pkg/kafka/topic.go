package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func EnsureTopicExists(bootstrapServers string, topic string, numPartitions int, replicationFactor int) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kafka admin client: %v", err)
	}
	defer adminClient.Close()

	metadata, err := adminClient.GetMetadata(&topic, false, 5000)
	if err != nil {
		return fmt.Errorf("failed to get Kafka metadata: %v", err)
	}

	if _, exists := metadata.Topics[topic]; !exists {
		topicConfig := []kafka.TopicSpecification{
			{
				Topic:             topic,
				NumPartitions:     numPartitions,
				ReplicationFactor: replicationFactor,
			},
		}

		ctx := context.Background()
		if _, err := adminClient.CreateTopics(ctx, topicConfig); err != nil {
			return fmt.Errorf("failed to create Kafka topic: %v", err)
		}
	}

	return nil
}
