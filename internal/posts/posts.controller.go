package posts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostController interface {
	RegisterRoutes(router *chi.Mux)
	CreatePost(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByUserId(w http.ResponseWriter, r *http.Request)
}

type postController struct {
	service PostService
}

func (c *postController) RegisterRoutes(router *chi.Mux) {
	router.Route("/posts", func(r chi.Router) {
		r.Get("/", c.FindAll)
		r.Get("/by-user/{userId}", c.FindByUserId)
		r.Post("/create", c.CreatePost)
	})
}

func (c *postController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if err := c.service.CreatePost(r.Context(), &post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (c *postController) FindAll(w http.ResponseWriter, r *http.Request) {
	posts, err := c.service.GetAllPosts(r.Context())

	if err != nil {
		http.Error(w, "error fetching posts", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (c *postController) FindByUserId(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(idParam)

	if err != nil {
		http.Error(w, "invalid userid", http.StatusBadRequest)
		return
	}

	posts, err := c.service.GetPostsByUserId(r.Context(), userId)

	if err != nil {
		http.Error(w, "error fetching posts", http.StatusBadGateway)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func NewPostController(service PostService) PostController {
	return &postController{service: service}
}
