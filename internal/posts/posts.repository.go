package posts

import (
	"context"
	"errors"
	"fmt"
	"go-back/internal/entities"
	"go-back/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	CreatePost(ctx context.Context, dto *CreatePostDto) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	GetPostsByUserId(ctx context.Context, userId uint64) ([]*models.Post, error)
	DeletePost(ctx context.Context, id uint64) error
}

type postRepository struct {
	db *sqlx.DB
}

func (r *postRepository) CreatePost(ctx context.Context, dto *CreatePostDto) (*models.Post, error) {
	query := `
		INSERT INTO posts (content, title, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, content, title, user_id, created_at
	`

	post := &entities.Post{}

	if err := r.db.GetContext(ctx, post, query, dto.Content, dto.Title, dto.UserID); err != nil {
		return nil, fmt.Errorf("error creating post %w", err)
	}

	return post.ToModel(), nil
}

func (r *postRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	query := `SELECT id, content, title, user_id, created_at FROM posts`

	posts := make([]*entities.Post, 0)
	if err := r.db.SelectContext(ctx, &posts, query); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get posts")
	}

	modelPosts := make([]*models.Post, 0)

	for _, post := range posts {
		modelPosts = append(modelPosts, post.ToModel())
	}

	return modelPosts, nil
}

func (r *postRepository) GetPostsByUserId(ctx context.Context, userId uint64) ([]*models.Post, error) {
	query := `SELECT id, content, title, user_id, created_at FROM posts WHERE user_id = $1`

	posts := make([]*entities.Post, 0)
	if err := r.db.SelectContext(ctx, posts, query, userId); err != nil {
		return nil, fmt.Errorf("failed to get posts")
	}

	modelPosts := make([]*models.Post, 0)
	for _, post := range posts {
		modelPosts = append(modelPosts, post.ToModel())
	}

	return modelPosts, nil
}

func (r *postRepository) DeletePost(ctx context.Context, id uint64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting post")
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		return errors.New(PostNotFound)
	}

	return nil
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}
