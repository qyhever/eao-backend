package repository

import (
	"eao/internal/model"
)

type PostRepository interface {
	GetPostList(query *model.PostListQuery) ([]model.Post, int, error)
	GetPostByID(id string) (*model.Post, error)
	CreatePost(param *model.CreatePostRequest) (string, error)
	UpdatePost(id string, param *model.UpdatePostRequest) error
	DeletePost(id string) error
}
