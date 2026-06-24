package service

import (
	"eao/internal/model"
	"eao/internal/repository"
)

type DiscService struct {
	repo repository.DiscRepository
}

func NewDiscService(repo repository.DiscRepository) *DiscService {
	return &DiscService{repo: repo}
}

func (s *DiscService) GetDiscList(query *model.DiscListQuery) ([]model.Disc, int, error) {
	return s.repo.GetDiscList(query)
}
