package repository

import "eao/internal/model"

type DiscRepository interface {
	GetDiscList(query *model.DiscListQuery) ([]model.Disc, int, error)
}
