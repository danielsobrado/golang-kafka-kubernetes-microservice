package config

import (
	"fmt"

	"github.com/magiconair/properties"
)

// Config represents the application configuration
type Config struct {
	ServerPort            int
	DatabaseURL           string
	KafkaBootstrapServers string
	KafkaTopic            string
	KafkaConsumerGroupID  string

	DatabaseQueries struct {
		GetAllUsers string
		GetUserByID string
		CreateOrder string
	}
}

// LoadConfig loads the configuration from the application.properties file
func LoadConfig(filePath string) (*Config, error) {
	props, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file: %v", err)
	}

	serverPort := props.GetInt("server.port", 0)
	if serverPort == 0 {
		return nil, fmt.Errorf("missing or invalid server port property")
	}

	dbURL := props.GetString("database.url", "")
	if dbURL == "" {
		return nil, fmt.Errorf("missing database URL property")
	}

	kafkaBootstrapServers := props.GetString("kafka.bootstrap.servers", "")
	if kafkaBootstrapServers == "" {
		return nil, fmt.Errorf("missing Kafka bootstrap servers property")
	}

	kafkaTopic := props.GetString("kafka.topic", "")
	if kafkaTopic == "" {
		return nil, fmt.Errorf("missing Kafka topic property")
	}

	kafkaConsumerGroupID := props.GetString("kafka.consumer.group.id", "")
	if kafkaConsumerGroupID == "" {
		return nil, fmt.Errorf("missing Kafka consumer group ID property")
	}

	return &Config{
		ServerPort:            serverPort,
		DatabaseURL:           dbURL,
		KafkaBootstrapServers: kafkaBootstrapServers,
		KafkaTopic:            kafkaTopic,
		KafkaConsumerGroupID:  kafkaConsumerGroupID,
	}, nil
}
