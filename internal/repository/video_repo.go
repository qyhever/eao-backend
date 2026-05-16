package repository

import "eao/internal/model"

type VideoRepository interface {
	GetVideoList() ([]model.VideoConfig, error)
}
