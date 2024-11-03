package users

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type UsersModule struct {
	Repo       UserRepository
	Service    UserService
	Controller UserController
}

func (m *UsersModule) Register(db *sql.DB, router *chi.Mux) {
	m.Repo = NewUserRepository(db)
	m.Service = NewUserService(m.Repo)
	m.Controller = NewUserController(m.Service)

	m.Controller.RegisterRoutes(router) 
}

func NewUsersModule() *UsersModule {
	return &UsersModule{}
}
