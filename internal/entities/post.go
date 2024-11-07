package entities

import "go-back/internal/models"

type Post struct {
	ID        uint64 `db:"id"`
	Content   string `db:"content"`
	Title     string `db:"title"`
	UserID    uint64 `db:"user_id"`
	CreatedAt string `db:"created_at"`
}

func (p *Post) ToModel() *models.Post {
	return &models.Post{
		ID:        p.ID,
		Content:   p.Content,
		Title:     p.Title,
		UserID:    p.UserID,
		CreatedAt: p.CreatedAt,
	}
}
