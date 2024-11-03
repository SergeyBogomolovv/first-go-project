package posts

import (
	"context"
	"errors"
	"fmt"
	"go-back/internal/users"
)

type PostService interface {
	CreatePost(ctx context.Context, post *Post) error
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostsByUserId(ctx context.Context, userId int) ([]*Post, error)
}

type postsService struct {
	postRepo PostRepository
	userService users.UserService
}

func (s *postsService) CreatePost(ctx context.Context, post *Post) error {
	switch {
	case post.Content == "":
		return errors.New("invalid content")
	case post.Title == "":
		return errors.New("invalid title")		
	}
	
	if _, err := s.userService.GetUserByID(ctx, post.UserID); err != nil {
		return fmt.Errorf("post author not found, id is %v", post.UserID)
	}

	if err := s.postRepo.CreatePost(ctx, post); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to create post")
	}

	return nil
}

func (s *postsService) GetAllPosts(ctx context.Context) ([]*Post, error) {
	return s.postRepo.GetAllPosts(ctx)
}

func (s *postsService) GetPostsByUserId(ctx context.Context, userId int) ([]*Post, error) {
	return s.postRepo.GetPostsByUserId(ctx, userId)
}

func NewPostService(postRepo PostRepository, userService users.UserService) PostService {
	return &postsService{postRepo: postRepo, userService: userService}
}