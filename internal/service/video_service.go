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
		if strings.HasPrefix(list[i].FileName, "http://") || strings.HasPrefix(list[i].FileName, "https://") {
			continue
		}
		list[i].FileName = baseURL + "/" + strings.TrimLeft(list[i].FileName, "/")
	}

	return list, nil
}
