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
