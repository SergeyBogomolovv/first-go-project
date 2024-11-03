package users

import (
	"context"
	"errors"
	"log"
)

type userService struct {
	repo UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, id int) error
}

func (s *userService) CreateUser(ctx context.Context, user *User) error {
	switch {
	case user.Email == "":
		return errors.New("invalid email")
	case len(user.Password) < 6:
		return errors.New("password must be longer than 6 chars")
	case user.Username == "":
		return errors.New("invalid username")
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*User, error) {
	log.Println("Getting user")
	user, err :=  s.repo.GetOne(ctx, id)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*User, error) {
	return s.repo.GetMany(ctx)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	log.Printf("User with id %v deleted", id)
	return nil
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}