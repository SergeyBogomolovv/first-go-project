package users

import (
	"encoding/json"
	"fmt"
	response "go-back/pkg/http"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserController interface {
	RegisterRoutes(mux *chi.Mux)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	service  UserService
	validate *validator.Validate
}

func (c *userController) RegisterRoutes(mux *chi.Mux) {
	mux.Route("/users", func(r chi.Router) {
		r.Get("/{id}", c.GetUserByID)
		r.Get("/", c.GetAllUsers)
		r.Post("/create", c.CreateUser)
		r.Delete("/{id}", c.DeleteUser)
	})
}

// GetUserByID godoc
//	@Summary	Get user details
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"User ID"
//	@Success	200	{object}	models.User
//	@Failure	400	{object}	response.ErrorResponse
//	@Router		/users/{id} [get]
func (c *userController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(r.Context(), uint64(id))
	if err != nil {
		response.SendError(w, "User not found", http.StatusNotFound)
		return
	}

	response.SendJSON(w, user, http.StatusOK)
}

// CreateUser godoc
//	@Summary	Create new user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		post	body		CreateUserDto	true	"Create new user"
//	@Success	201		{object}	models.User
//	@Failure	400		{object}	response.ErrorResponse
//	@Router		/users/create [post]
func (c *userController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var dto CreateUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := c.validate.Struct(&dto); err != nil {
		response.SendError(w, fmt.Sprintf("Validation failed: %s", err), http.StatusBadRequest)
		return
	}

	user, err := c.service.CreateUser(r.Context(), &dto)

	if err != nil {
		if err.Error() == UserExists {
			response.SendError(w, "User already exists", http.StatusConflict)
		} else {
			response.SendError(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	response.SendJSON(w, user, http.StatusCreated)
}

// GetAllUsers godoc
//	@Summary	Get all users
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	models.User
//	@Router		/users [get]
func (c *userController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.GetAllUsers(r.Context())
	if err != nil {
		response.SendError(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	response.SendJSON(w, users, http.StatusOK)
}

// DeleteUser godoc
//	@Summary	Delete user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"User ID"
//	@Failure	400	{object}	response.ErrorResponse
//	@Router		/{id} [delete]
func (c *userController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(r.Context(), uint64(id)); err != nil {
		if err.Error() == UserNotFound {
			response.SendError(w, "User not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "error deleting user", http.StatusInternalServerError)
		return
	}

	response.SendMessage(w, "User deleted successfully", http.StatusOK)
}

func NewUserController(service UserService) UserController {
	return &userController{service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}
