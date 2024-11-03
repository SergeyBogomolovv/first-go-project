package posts

import (
	"database/sql"
	"go-back/internal/users"

	"github.com/go-chi/chi/v5"
)

type PostsModule struct {
	Repo PostRepository
	Service PostService
	Controller PostController
	userService users.UserService
}

func (m *PostsModule) InjectUserService(s users.UserService) {
	m.userService = s
}

func (m *PostsModule) Register(db *sql.DB, router *chi.Mux) {
	m.Repo = NewPostRepository(db)
	m.Service = NewPostService(m.Repo, m.userService)
	m.Controller = NewPostController(m.Service)

	m.Controller.RegisterRoutes(router) 
}

func NewPostsModule() *PostsModule {
	return &PostsModule{}
}
