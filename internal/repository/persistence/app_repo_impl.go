package persistence

import (
	"eao/internal/model"
	"eao/internal/repository"
)

type AppRepositoryImpl struct{}

func NewAppRepository() repository.AppRepository {
	return &AppRepositoryImpl{}
}

func (r *AppRepositoryImpl) GetHelloInfo(req *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error) {
	res := &model.GetHelloInfoResponse{
		Name: req.Name,
	}
	return res, nil
}
