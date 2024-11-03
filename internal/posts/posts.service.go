package posts

type PostService interface {}

type postsService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) PostService {
	return &postsService{repo: repo}
}