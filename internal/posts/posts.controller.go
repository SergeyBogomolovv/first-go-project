package posts

import "github.com/go-chi/chi/v5"

type PostController interface {
	RegisterRoutes(router *chi.Mux)
}

type postController struct {
	service PostService
}

func (c *postController) RegisterRoutes(router *chi.Mux) {
	router.Route("/posts", func(r chi.Router) {})
}

func NewPostController(service PostService) PostController {
	return &postController{service: service}
}