package repository

import (
	"context"
	"errors"

	"eao/internal/model"
)

var ErrAdminUsernameAlreadyExists = errors.New("admin username already exists")

type AdminRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.Admin, error)
	FindByID(ctx context.Context, id int64) (*model.Admin, error)
	Upsert(ctx context.Context, admin model.Admin) error
	Create(ctx context.Context, admin model.Admin) (*model.Admin, error)
}
