package repository

import (
	"database/sql"
	"fmt"
	"golang-kafka-kubernetes-microservice/pkg/config"
	"golang-kafka-kubernetes-microservice/pkg/model"
)

// Repository represents the data access layer
type Repository struct {
	db     *sql.DB
	config *config.Config
}

// NewRepository creates a new instance of the Repository
func NewRepository(db *sql.DB, config *config.Config) *Repository {
	return &Repository{
		db:     db,
		config: config,
	}
}

// GetAllUsers retrieves all users from the database
func (r *Repository) GetAllUsers() ([]*model.User, error) {
	query := r.config.DatabaseQueries.GetAllUsers
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return users, nil
}

// GetUserByID retrieves a user by ID from the database
func (r *Repository) GetUserByID(userID int) (*model.User, error) {
	query := r.config.DatabaseQueries.GetUserByID
	row := r.db.QueryRow(query, userID)

	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %v", err)
	}

	return &user, nil
}

// CreateOrder creates a new order in the database
func (r *Repository) CreateOrder(order *model.Order) error {
	query := r.config.DatabaseQueries.CreateOrder
	err := r.db.QueryRow(query, order.UserID, order.Total, order.Status, order.CreatedAt).Scan(&order.ID)
	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}

	return nil
}
