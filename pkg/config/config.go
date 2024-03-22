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
}

// LoadConfig loads the configuration from the application.properties file
func LoadConfig(filePath string) (*Config, error) {
	props, err := properties.LoadFile(filePath, properties.UTF8)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file: %v", err)
	}

	serverPort, err := props.GetInt("server.port")
	if err != nil {
		return nil, fmt.Errorf("invalid server port: %v", err)
	}

	dbURL, err := props.GetString("database.url")
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %v", err)
	}

	kafkaBootstrapServers, err := props.GetString("kafka.bootstrap.servers")
	if err != nil {
		return nil, fmt.Errorf("invalid Kafka bootstrap servers: %v", err)
	}

	kafkaTopic, err := props.GetString("kafka.topic")
	if err != nil {
		return nil, fmt.Errorf("invalid Kafka topic: %v", err)
	}

	return &Config{
		ServerPort:            serverPort,
		DatabaseURL:           dbURL,
		KafkaBootstrapServers: kafkaBootstrapServers,
		KafkaTopic:            kafkaTopic,
	}, nil
}
