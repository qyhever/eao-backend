package service

import (
	"eao/internal/config"
	"eao/internal/model"
	"eao/internal/repository"
	"strings"
)

type VideoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) GetVideoList() ([]model.VideoConfig, error) {
	list, err := s.repo.GetVideoList()
	if err != nil {
		return nil, err
	}

	cfg := config.GetConfig()
	if cfg == nil || strings.TrimSpace(cfg.PublicBaseURL) == "" {
		return list, nil
	}

	baseURL := strings.TrimRight(cfg.PublicBaseURL, "/")
	for i := range list {
		list[i].FileName = withPublicBaseURL(baseURL, list[i].FileName)
		list[i].Cover = withPublicBaseURL(baseURL, list[i].Cover)
	}

	return list, nil
}

func withPublicBaseURL(baseURL, name string) string {
	if name == "" || strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		return name
	}
	return baseURL + "/" + strings.TrimLeft(name, "/")
}
