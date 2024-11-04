package users

import (
	"context"
	"fmt"
	"go-back/internal/entities"
	"go-back/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, dto *CreateUserDto) (*models.User, error)
	GetOne(ctx context.Context, id uint64) (*models.User, error)
	GetMany(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	CheckUserExists(ctx context.Context, dto *UserExistsDto) (bool, error)
}

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) Create(ctx context.Context, dto *CreateUserDto) (*models.User, error) {
	query := `
		INSERT INTO users (username, password, email) 
		VALUES ($1, $2, $3) 
		RETURNING id, username, email, created_at
	`
	user := &entities.User{}

	if err := r.db.GetContext(ctx, user, query, dto.Username, dto.Password, dto.Email); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return user.ToModel(), nil
}

func (r *userRepository) CheckUserExists(ctx context.Context, dto *UserExistsDto) (bool, error) {
	query := `
		SELECT CASE WHEN EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2) THEN TRUE ELSE FALSE END AS record_exists		
	`

	var isUserExists bool
	err := r.db.GetContext(ctx, &isUserExists, query, dto.Username, dto.Email)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %w", err)
	}

	return isUserExists, nil
}

func (r *userRepository) GetOne(ctx context.Context, id uint64) (*models.User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	user := &entities.User{}

	if err := r.db.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}

	return user.ToModel(), nil
}

func (r *userRepository) GetMany(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, email, created_at FROM users`
	users := make([]*entities.User, 0)

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	modelUsers := make([]*models.User, 0)
	for _, user := range users {
		modelUsers = append(modelUsers, user.ToModel())
	}

	return modelUsers, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := `DELETE FROM users WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user")
	}

	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}
