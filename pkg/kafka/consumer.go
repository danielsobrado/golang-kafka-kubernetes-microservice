package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// See: https://github.com/confluentinc/confluent-kafka-go/issues/1014
type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

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

func (c *Consumer) Close() {
	c.consumer.Close()
}
