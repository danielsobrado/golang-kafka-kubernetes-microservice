package service

import (
	"encoding/json"
	"fmt"
	"golang-kafka-kubernetes-microservice/pkg/model"
	"golang-kafka-kubernetes-microservice/pkg/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllUsers() ([]*model.User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetUserByID(userID int) (*model.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *Service) CreateOrder(order *model.Order) error {
	user, err := s.repo.GetUserByID(order.UserID)
	if err != nil {
		return err
	}

	if user.Balance < order.Total {
		return fmt.Errorf("insufficient funds")
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return err
	}

	user.Balance -= order.Total
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *Service) ProcessOrder(order *model.Order) error {
	producer, err := kafka.NewProducer(kafkaBootstrapServers)
	if err != nil {
		return err
	}
	defer producer.Close()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	if err := producer.Produce(kafkaTopic, orderJSON); err != nil {
		return err
	}

	return nil
}
