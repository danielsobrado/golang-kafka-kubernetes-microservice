package service

import (
	"encoding/json"
	"fmt"
	"golang-kafka-kubernetes-microservice/pkg/model"
	"golang-kafka-kubernetes-microservice/pkg/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repo                  *repository.Repository
	kafkaBootstrapServers string
	kafkaTopic            string
}

func NewService(repo *repository.Repository, kafkaBootstrapServers string, kafkaTopic string) *Service {
	return &Service{
		repo:                  repo,
		kafkaBootstrapServers: kafkaBootstrapServers,
		kafkaTopic:            kafkaTopic,
	}
}

func (s *Service) GetAllUsers() ([]*model.User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetUserByID(userID int) (*model.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *Service) CreateOrder(order *model.Order) (*model.Order, error) {
	user, err := s.repo.GetUserByID(order.UserID)
	if err != nil {
		return nil, err
	}

	if user.Balance < order.Total {
		return nil, fmt.Errorf("insufficient funds")
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	user.Balance -= order.Total
	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	if err := s.ProcessOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) ProcessOrder(order *model.Order) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": s.kafkaBootstrapServers,
	})
	if err != nil {
		return err
	}
	defer producer.Close()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.kafkaTopic, Partition: kafka.PartitionAny},
		Value:          orderJSON,
	}
	if err := producer.Produce(message, nil); err != nil {
		return err
	}

	return nil
}
