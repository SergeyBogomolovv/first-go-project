package users

import (
	"encoding/json"
	"go-back/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController interface {
	RegisterRoutes(mux *chi.Mux)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	service UserService
}

func (c *userController) RegisterRoutes(mux *chi.Mux) {
	mux.Route("/users", func(r chi.Router) {
		r.Get("/{id}", c.GetUserByID)
		r.Get("/", c.GetAllUsers)
		r.Post("/create", c.CreateUser)
		r.Delete("/{id}", c.DeleteUser)
	})
}

func (c *userController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(r.Context(), uint64(id))
	if err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *userController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto CreateUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := c.service.CreateUser(r.Context(), &dto)

	if err != nil {
		if err.Error() == UserExists {
			utils.SendErrorResponse(w, "User already exists", http.StatusConflict)
		} else {
			utils.SendErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *userController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAllUsers(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(r.Context(), uint64(id)); err != nil {
		if err.Error() == UserNotFound {
			utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}
		utils.SendErrorResponse(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	utils.SendMessageReponse(w, "User deleted successfully", http.StatusOK)
}

func NewUserController(service UserService) UserController {
	return &userController{service: service}
}
