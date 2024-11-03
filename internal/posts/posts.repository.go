package posts

import (
	"context"
	"database/sql"
	"fmt"
)

type Post struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type PostRepository interface {
	CreatePost(ctx context.Context, post *Post) error
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostsByUserId(ctx context.Context, userId int) ([]*Post, error)
}

type postRepository struct {
	db *sql.DB
}

func (r *postRepository) CreatePost(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID).Scan(&post.ID, &post.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating post %w", err)
	}
	return nil
}

func (r *postRepository) GetAllPosts(ctx context.Context) ([]*Post, error) {
	query := `SELECT id, content, title, user_id, created_at FROM posts`
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to get posts")
	}

	return getPostsArray(rows)
}

func (r *postRepository) GetPostsByUserId(ctx context.Context, userId int) ([]*Post, error) {
	query := `SELECT id, content, title, user_id, created_at FROM posts WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to get posts")
	}

	return getPostsArray(rows)
}

func getPostsArray(rows *sql.Rows) ([]*Post, error) {
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := rows.Scan(&post.ID, &post.Content, &post.Title, &post.UserID, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning posts")
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading posts: %w", err)
	}

	return posts, nil
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}