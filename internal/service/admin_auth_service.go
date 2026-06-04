package service

import (
	"context"
	"errors"
	"strings"

	"eao/internal/model"
	jwtpkg "eao/internal/pkg/jwt"
	"eao/internal/pkg/password"
	"eao/internal/repository"
)

type AdminAuthService struct {
	adminRepo repository.AdminAccountRepository
}

func NewAdminAuthService(
	adminRepo repository.AdminAccountRepository,
) *AdminAuthService {
	return &AdminAuthService{
		adminRepo: adminRepo,
	}
}

var (
	ErrInvalidAdminCredentials = errors.New("账号或密码错误")
	ErrInvalidRefreshToken     = errors.New("refresh token 无效")
	ErrUserNotFound            = errors.New("用户不存在")
)

func (s *AdminAuthService) AdminLogin(ctx context.Context, req model.AdminLoginRequest) (*model.AdminLoginResponse, error) {
	if s.adminRepo == nil {
		return nil, ErrAdminNotFound
	}

	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, ErrAdminUsernameRequired
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, ErrAdminPasswordRequired
	}

	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrInvalidAdminCredentials
	}
	if admin.Status != "" && admin.Status != "active" {
		return nil, ErrInvalidAdminCredentials
	}
	if err := password.Compare(admin.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidAdminCredentials
	}

	accessToken, refreshToken, err := jwtpkg.GenToken(uint64(admin.ID))
	if err != nil {
		return nil, err
	}
	return &model.AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
