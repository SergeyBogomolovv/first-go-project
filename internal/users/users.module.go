package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type UsersModule struct {
	repo       UserRepository
	service    UserService
	controller *userController
}

func (m *UsersModule) Register(db *sql.DB, router *chi.Mux) {
	m.repo = NewUserRepository(db)
	m.service = NewUserService(m.repo)
	m.controller = NewUserController(m.service)

	m.controller.RegisterRoutes(router) 
}

func NewUsersModule() *UsersModule {
	return &UsersModule{}
}
