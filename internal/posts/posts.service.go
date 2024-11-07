package posts

import (
	"context"
	"errors"
	"fmt"
	"go-back/internal/models"
	"go-back/internal/users"
)

type PostService interface {
	CreatePost(ctx context.Context, dto *CreatePostDto) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	GetPostsByUserId(ctx context.Context, userId uint64) ([]*models.Post, error)
	DeletePost(ctx context.Context, id uint64) error
}

type postsService struct {
	postRepo    PostRepository
	userService users.UserService
}

func (s *postsService) CreatePost(ctx context.Context, dto *CreatePostDto) (*models.Post, error) {
	if _, err := s.userService.GetUserByID(ctx, dto.UserID); err != nil {
		return nil, errors.New(AuthorNotFound)
	}

	post, err := s.postRepo.CreatePost(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to create post")
	}

	return post, nil
}

func (s *postsService) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	return s.postRepo.GetAllPosts(ctx)
}

func (s *postsService) GetPostsByUserId(ctx context.Context, userId uint64) ([]*models.Post, error) {
	return s.postRepo.GetPostsByUserId(ctx, userId)
}

func (s *postsService) DeletePost(ctx context.Context, id uint64) error {
	return s.postRepo.DeletePost(ctx, id)
}

func NewPostService(postRepo PostRepository, userService users.UserService) PostService {
	return &postsService{postRepo: postRepo, userService: userService}
}
