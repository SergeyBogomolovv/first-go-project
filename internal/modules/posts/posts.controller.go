package posts

import (
	"encoding/json"
	"fmt"
	response "go-back/pkg/http"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type PostController interface {
	RegisterRoutes(router *chi.Mux)
	CreatePost(w http.ResponseWriter, r *http.Request)
	FindAllPosts(w http.ResponseWriter, r *http.Request)
	FindByUserId(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
}

type postController struct {
	service  PostService
	validate *validator.Validate
}

func (c *postController) RegisterRoutes(router *chi.Mux) {
	router.Route("/posts", func(r chi.Router) {
		r.Get("/", c.FindAllPosts)
		r.Get("/by-user/{userId}", c.FindByUserId)
		r.Post("/create", c.CreatePost)
		r.Delete("/delete/{id}", c.DeletePost)
	})
}

// CreatePost godoc
//	@Summary	Create new post
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Param		post	body		CreatePostDto	true	"Create new post"
//	@Success	201		{object}	models.Post
//	@Failure	400		{object}	response.ErrorResponse
//	@Router		/posts/create [post]
func (c *postController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var dto CreatePostDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := c.validate.Struct(&dto); err != nil {
		response.SendError(w, fmt.Sprintf("Validation failed: %s", err), http.StatusBadRequest)
		return
	}

	post, err := c.service.CreatePost(r.Context(), &dto)

	if err != nil {
		if err.Error() == AuthorNotFound {
			response.SendError(w, "Author not found", http.StatusNotFound)
			return
		}
		response.SendError(w, fmt.Sprintf("Error creating post: %s", err), http.StatusInternalServerError)
		return
	}

	response.SendJSON(w, post, http.StatusCreated)
}

// FindAllPosts godoc
//	@Summary	Get details of all posts
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		models.Post
//	@Failure	400	{object}	response.ErrorResponse
//	@Router		/posts [get]
func (c *postController) FindAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.service.GetAllPosts(r.Context())

	if err != nil {
		response.SendError(w, "error fetching posts", http.StatusBadGateway)
		return
	}
	response.SendJSON(w, posts, http.StatusOK)
}

// FindByUserId godoc
//	@Summary	Get details users posts
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Param		userId	path		int	true	"User ID"
//	@Success	200		{array}		models.Post
//	@Failure	400		{object}	response.ErrorResponse
//	@Router		/posts/by-user/{userId} [get]
func (c *postController) FindByUserId(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "userId")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.SendError(w, "invalid userid", http.StatusBadRequest)
		return
	}

	posts, err := c.service.GetPostsByUserId(r.Context(), uint64(userId))

	if err != nil {
		response.SendError(w, "error fetching posts", http.StatusBadGateway)
		return
	}
	response.SendJSON(w, posts, http.StatusOK)
}

// DeletePost godoc
//	@Summary	Get details users posts
//	@Tags		posts
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Post ID"
//	@Failure	400	{object}	response.ErrorResponse
//	@Router		/posts/delete/{id} [delete]
func (c *postController) DeletePost(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.SendError(w, "invalid post id", http.StatusBadRequest)
		return
	}

	if err := c.service.DeletePost(r.Context(), uint64(id)); err != nil {
		if err.Error() == PostNotFound {
			response.SendError(w, "Post not found", http.StatusNotFound)
			return
		}
		response.SendError(w, "error deleting post", http.StatusInternalServerError)
		return
	}
	response.SendMessage(w, "Post deleted successfully", http.StatusOK)
}

func NewPostController(service PostService) PostController {
	return &postController{service: service, validate: validator.New(validator.WithRequiredStructEnabled())}
}
