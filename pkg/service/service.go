package service

import (
	"encoding/json"
	"fmt"
	"golang-kafka-kubernetes-microservice/pkg/model"
	"golang-kafka-kubernetes-microservice/pkg/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Service represents the business logic layer
type Service struct {
	repo *repository.Repository
}

// NewService creates a new instance of the Service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllUsers retrieves all users
func (s *Service) GetAllUsers() ([]*model.User, error) {
	return s.repo.GetAllUsers()
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(userID int) (*model.User, error) {
	return s.repo.GetUserByID(userID)
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(order *model.Order) error {
	// Perform any necessary validations or business logic
	// For example, check if the user exists and has sufficient funds

	// Retrieve the user by ID
	user, err := s.repo.GetUserByID(order.UserID)
	if err != nil {
		return err
	}

	// Check if the user has sufficient funds
	if user.Balance < order.Total {
		return fmt.Errorf("insufficient funds")
	}

	// Create the order in the repository
	if err := s.repo.CreateOrder(order); err != nil {
		return err
	}

	// Update the user's balance
	user.Balance -= order.Total
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// ProcessOrder processes an order (e.g., sends it to Kafka)
func (s *Service) ProcessOrder(order *model.Order) error {
	// Perform any necessary processing logic
	// For example, send the order to a Kafka topic

	// Create a Kafka producer
	producer, err := kafka.NewProducer(kafkaBootstrapServers)
	if err != nil {
		return err
	}
	defer producer.Close()

	// Serialize the order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Send the order to the Kafka topic
	if err := producer.Produce(kafkaTopic, orderJSON); err != nil {
		return err
	}

	return nil
}
