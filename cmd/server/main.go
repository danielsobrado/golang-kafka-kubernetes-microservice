package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-kafka-kubernetes-microservice/pkg/config"
	"golang-kafka-kubernetes-microservice/pkg/db"
	"golang-kafka-kubernetes-microservice/pkg/handler"
	"golang-kafka-kubernetes-microservice/pkg/kafka"
	"golang-kafka-kubernetes-microservice/pkg/repository"
	"golang-kafka-kubernetes-microservice/pkg/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Connect to PostgreSQL database
	db, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure Kafka topic exists
	if err := kafka.EnsureTopicExists(cfg.KafkaBootstrapServers, cfg.KafkaTopic); err != nil {
		logger.Fatalf("Failed to ensure Kafka topic exists: %v", err)
	}

	// Create repository instance
	repo := repository.NewRepository(db)

	// Create service instance
	svc := service.NewService(repo)

	// Create Kafka consumer instance
	consumer, err := kafka.NewConsumer(cfg.KafkaBootstrapServers, cfg.KafkaTopic, svc)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create HTTP handler instance
	httpHandler := handler.NewHandler(svc)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      httpHandler,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
	}

	// Start the server
	go func() {
		logger.Printf("Starting server on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start Kafka consumer
	go func() {
		logger.Printf("Starting Kafka consumer")
		if err := consumer.Start(); err != nil {
			logger.Fatalf("Failed to start Kafka consumer: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown Kafka consumer
	consumer.Stop()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown server: %v", err)
	}

	logger.Println("Server stopped")
}