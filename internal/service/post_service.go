package service

import (
	"eao/internal/model"
	"eao/internal/repository"
)

type PostService struct {
	// 依赖接口，而不是具体结构体
	repo repository.PostRepository
}

// 构造函数注入依赖
func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPostList(query *model.PostListQuery) ([]model.Post, int, error) {
	return s.repo.GetPostList(query)
}

func (s *PostService) GetPostByID(id string) (*model.Post, error) {
	return s.repo.GetPostByID(id)
}

func (s *PostService) CreatePost(param *model.CreatePostRequest) (string, error) {
	return s.repo.CreatePost(param)
}

func (s *PostService) UpdatePost(id string, param *model.UpdatePostRequest) error {
	return s.repo.UpdatePost(id, param)
}

func (s *PostService) DeletePost(id string) error {
	return s.repo.DeletePost(id)
}
