package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userController struct {
	service UserService
}

func (c *userController) RegisterRoutes(mux *chi.Mux) {
	mux.Route("/users", func(r chi.Router) {
		r.Get("/{id}", c.getUserByID)
		r.Get("/", c.getAllUsers)
		r.Post("/create", c.createUser)
		r.Delete("/{id}", c.deleteUser)
	})
}

func (c *userController) getUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *userController) createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := c.service.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR: %v", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *userController) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *userController) deleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	message := struct{
		Message string `json:"message"`
	}{
		Message: "User deleted successfully",
	}
	json.NewEncoder(w).Encode(message)
}

func NewUserController(service UserService) *userController {
	return &userController{service: service}
}