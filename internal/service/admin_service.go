package service

import (
	"context"
	"errors"
	"strings"

	"eao/internal/model"
	"eao/internal/pkg/password"
	"eao/internal/repository"
)

var (
	ErrAdminUsernameRequired      = errors.New("管理员账号不能为空")
	ErrAdminPasswordRequired      = errors.New("管理员密码不能为空")
	ErrAdminUsernameAlreadyExists = errors.New("管理员账号已存在")
)

type AdminService struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

func (s *AdminService) CreateAdmin(ctx context.Context, req model.CreateAdminRequest) (*model.AdminProfileResponse, error) {

	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, ErrAdminUsernameRequired
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, ErrAdminPasswordRequired
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = username
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	admin, err := s.adminRepo.Create(ctx, model.Admin{
		Username:     username,
		PasswordHash: hash,
		Name:         name,
		Status:       "active",
	})
	if err != nil {
		if errors.Is(err, repository.ErrAdminUsernameAlreadyExists) {
			return nil, ErrAdminUsernameAlreadyExists
		}
		return nil, err
	}

	return &model.AdminProfileResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Name:     admin.Name,
		Status:   admin.Status,
	}, nil
}
