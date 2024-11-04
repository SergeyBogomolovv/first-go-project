package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetOne(ctx context.Context, id int) (*User, error)
	GetMany(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	existingUserQuery := `
		SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2) THEN TRUE ELSE FALSE END AS record_exists		
	`
	var isUserExists bool
	r.db.QueryRowContext(ctx, existingUserQuery, user.Username, user.Email).Scan(&isUserExists)
	if isUserExists {
		return errors.New("user already exists")
	}

	query := `
		INSERT INTO users (username, password, email) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *userRepository) GetOne(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, username, email, created_at 
		FROM users 
		WHERE id = $1
	`
	user := &User{}

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetMany(ctx context.Context) ([]*User, error) {
	query := `SELECT id, username, email, created_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading users: %w", err)
	}

	return users, nil
}

func (r *userRepository) DeleteUser(ctx context.Context,id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user")
	}
	
	return nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}