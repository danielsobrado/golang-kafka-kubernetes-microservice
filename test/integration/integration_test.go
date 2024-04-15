package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-kafka-kubernetes-microservice/pkg/config"
	"golang-kafka-kubernetes-microservice/pkg/db"
	"golang-kafka-kubernetes-microservice/pkg/handler"
	"golang-kafka-kubernetes-microservice/pkg/model"
	"golang-kafka-kubernetes-microservice/pkg/repository"
	"golang-kafka-kubernetes-microservice/pkg/service"
)

func TestGetAllUsers(t *testing.T) {
	// Load the configuration
	cfg, err := config.LoadConfig("../../application.properties")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the test database
	db, err := db.NewPostgresDB(cfg.GetString("test.database.url"))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Create the repository and service instances
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	// Create a test HTTP server
	server := httptest.NewServer(handler.NewHandler(svc))
	defer server.Close()

	// Send a GET request to the /users endpoint
	resp, err := http.Get(server.URL + "/users")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code. Got %d, expected %d", resp.StatusCode, http.StatusOK)
	}

	// Parse the response body
	var users []*model.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Assert the response body
	expectedUsers := []*model.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	if len(users) != len(expectedUsers) {
		t.Errorf("Unexpected number of users. Got %d, expected %d", len(users), len(expectedUsers))
	}
	for i, user := range users {
		if user.ID != expectedUsers[i].ID || user.Name != expectedUsers[i].Name || user.Email != expectedUsers[i].Email {
			t.Errorf("Unexpected user data. Got %+v, expected %+v", user, expectedUsers[i])
		}
	}
}
