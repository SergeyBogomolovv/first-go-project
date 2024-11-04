package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type UsersModule struct {
	Repo       UserRepository
	Service    UserService
	Controller UserController
}

func (m *UsersModule) Register(db *sqlx.DB, router *chi.Mux) {
	m.Repo = NewUserRepository(db)
	m.Service = NewUserService(m.Repo)
	m.Controller = NewUserController(m.Service)

	m.Controller.RegisterRoutes(router) 
}

func NewUsersModule() *UsersModule {
	return &UsersModule{}
}
