package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang-kafka-kubernetes-microservice/pkg/model"
	"golang-kafka-kubernetes-microservice/pkg/service"

	"github.com/go-chi/chi"
)

// Handler represents the HTTP handler for the microservice
type Handler struct {
	service *service.Service
	router  *chi.Mux
}

// NewHandler creates a new instance of the Handler
func NewHandler(service *service.Service) *Handler {
	h := &Handler{
		service: service,
		router:  chi.NewRouter(),
	}

	h.router.Get("/users", h.GetUsers)
	h.router.Get("/users/{id}", h.GetUserByID)
	h.router.Post("/orders", h.CreateOrder)

	return h
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// GetUsers handles the request to get all users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, users)
}

// GetUserByID handles the request to get a user by ID
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// CreateOrder handles the request to create a new order
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.CreateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusCreated, createdOrder)
}

// respondJSON sends a JSON response with the specified status code and payload
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
