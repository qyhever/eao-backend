package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"eao/internal/repository"
)

var ErrFileInvalidParam = errors.New("请求参数错误")

type FileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) List(ctx context.Context) ([]byte, error) {
	return s.repo.List(ctx)
}

func (s *FileService) ListByDir(ctx context.Context, dirName string) ([]byte, error) {
	if strings.TrimSpace(dirName) == "" {
		return nil, ErrFileInvalidParam
	}
	return s.repo.ListByDir(ctx, dirName)
}

func (s *FileService) Upload(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, dirName string) ([]byte, error) {
	if file == nil || fileHeader == nil || strings.TrimSpace(dirName) == "" {
		return nil, ErrFileInvalidParam
	}
	return s.repo.Upload(ctx, file, fileHeader, dirName)
}
