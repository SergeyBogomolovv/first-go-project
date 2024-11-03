package store

import (
	"context"
	"database/sql"
)

type Post struct {
	ID int64 `json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int64 `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID).Scan(&post.ID, &post.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}