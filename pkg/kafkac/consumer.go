package kafkac

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

// NewConsumer creates a new instance of the Kafka consumer
func NewConsumer(bootstrapServers string, groupID string, topic string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

// Consume starts consuming messages from the Kafka topic
func (c *Consumer) Consume(handler func(message *kafka.Message)) error {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("failed to subscribe to Kafka topic: %v", err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			return fmt.Errorf("failed to read message from Kafka topic: %v", err)
		}

		handler(msg)
	}
}

// Close closes the Kafka consumer
func (c *Consumer) Close() {
	c.consumer.Close()
}
