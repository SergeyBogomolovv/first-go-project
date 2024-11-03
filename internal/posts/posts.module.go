package posts

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type PostsModule struct {
	repo PostRepository
	service PostService
	controller PostController
}

func (m *PostsModule) Register(db *sql.DB, router *chi.Mux) {
	m.repo = NewPostRepository(db)
	m.service = NewPostService(m.repo)
	m.controller = NewPostController(m.service)

	m.controller.RegisterRoutes(router) 
}

func NewPostsModule() *PostsModule {
	return &PostsModule{}
}
