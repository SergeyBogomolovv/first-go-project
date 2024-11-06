package users

import (
	"context"
	"errors"
	"fmt"
	"go-back/internal/models"
)

type userService struct {
	repo UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, dto *CreateUserDto) (*models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

func (s *userService) CreateUser(ctx context.Context, dto *CreateUserDto) (*models.User, error) {
	isUserExists, err := s.repo.CheckUserExists(ctx, &UserExistsDto{Username: dto.Username, Email: dto.Email})
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists: %w", err)
	}
	if isUserExists {
		return nil, errors.New(UserExists)
	}

	return s.repo.Create(ctx, dto)
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetOne(ctx, id)

	if err != nil {
		return nil, errors.New(UserNotFound)
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.GetMany(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	return users, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}
