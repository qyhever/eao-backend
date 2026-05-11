package repository

import (
	"eao/internal/model"
)

type AppRepository interface {
	GetHelloInfo(param *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error)
}
