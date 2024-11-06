package posts

type CreatePostDto struct {
	Content string `json:"content" validate:"required"`
	Title   string `json:"title" validate:"required"`
	UserID  uint64 `json:"user_id" validate:"required"`
}

const (
	PostNotFound = "post not found"
)
