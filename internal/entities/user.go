package entities

import "go-back/internal/models"

type User struct {
	ID        uint64 `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	CreatedAt string `db:"created_at"`
}

func (u *User) ToModel() *models.User {
	return &models.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

type UserWithPost struct {
	ID            uint64 `db:"id"`
	Username      string `db:"username"`
	Email         string `db:"email"`
	CreatedAt     string `db:"created_at"`
	PostID        uint64 `db:"post_id"`
	Content       string `db:"post_content"`
	Title         string `db:"post_title"`
	PostCreatedAt string `db:"post_created_at"`
}
