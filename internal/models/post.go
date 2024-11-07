package models

type Post struct {
	ID        uint64 `json:"id"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	UserID    uint64 `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
