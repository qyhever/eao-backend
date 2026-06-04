package service

import (
	"context"
	"errors"
	"strings"

	"eao/internal/domain"
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
	ErrInvalidAdminCredentials = domain.ErrInvalidAdminCredentials
	ErrInvalidRefreshToken     = domain.ErrInvalidRefreshToken
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
		return nil, domain.ErrInvalidAdminCredentials
	}
	if admin.Status != "" && admin.Status != "active" {
		return nil, domain.ErrInvalidAdminCredentials
	}
	if err := password.Compare(admin.PasswordHash, req.Password); err != nil {
		return nil, domain.ErrInvalidAdminCredentials
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

func (s *AdminAuthService) AdminRefreshToken(ctx context.Context, req model.AdminRefreshTokenRequest) (*model.AdminLoginResponse, error) {
	if s.adminRepo == nil {
		return nil, ErrAdminNotFound
	}

	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, domain.ErrInvalidRefreshToken
	}

	claims, err := jwtpkg.ParseToken(refreshToken)
	if err != nil || !claims.IsRefreshToken() {
		return nil, domain.ErrInvalidRefreshToken
	}

	admin, err := s.adminRepo.FindByID(ctx, int64(claims.UserID))
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrUserNotFound
	}
	if admin.Status != "" && admin.Status != "active" {
		return nil, domain.ErrInvalidAdminCredentials
	}

	accessToken, newRefreshToken, err := jwtpkg.GenToken(uint64(admin.ID))
	if err != nil {
		return nil, err
	}
	return &model.AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
