package service

import (
	"context"
	"errors"
	"testing"

	"eao/internal/config"
	"eao/internal/model"
	jwtpkg "eao/internal/pkg/jwt"
	"eao/internal/repository/persistence"
)

func TestAdminAuthServiceAdminLoginSuccess(t *testing.T) {
	config.GlobalConfig = testAuthConfig()
	ctx := context.Background()
	repo := persistence.NewAdminAccountRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminAuthService(repo)

	result, err := svc.AdminLogin(ctx, model.AdminLoginRequest{
		Username: " admin ",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("expected tokens, got %+v", result)
	}

	claims, err := jwtpkg.ParseToken(result.AccessToken)
	if err != nil {
		t.Fatalf("parse access token failed: %v", err)
	}
	if claims.UserID != 1 || !claims.IsAccessToken() {
		t.Fatalf("unexpected access token claims: %+v", claims)
	}
}

func TestAdminAuthServiceAdminLoginInvalidPassword(t *testing.T) {
	config.GlobalConfig = testAuthConfig()
	ctx := context.Background()
	repo := persistence.NewAdminAccountRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminAuthService(repo)

	_, err := svc.AdminLogin(ctx, model.AdminLoginRequest{
		Username: "admin",
		Password: "bad-password",
	})
	if !errors.Is(err, ErrInvalidAdminCredentials) {
		t.Fatalf("expected ErrInvalidAdminCredentials, got %v", err)
	}
}

func TestAdminAuthServiceAdminLoginDisabledAdmin(t *testing.T) {
	config.GlobalConfig = testAuthConfig()
	ctx := context.Background()
	repo := persistence.NewAdminAccountRepository(nil)
	admin := testAdmin(1, "admin")
	admin.Status = "disabled"
	if err := repo.Upsert(ctx, admin); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminAuthService(repo)

	_, err := svc.AdminLogin(ctx, model.AdminLoginRequest{
		Username: "admin",
		Password: "password",
	})
	if !errors.Is(err, ErrInvalidAdminCredentials) {
		t.Fatalf("expected ErrInvalidAdminCredentials, got %v", err)
	}
}

func testAuthConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test-secret",
			AccessExpiresIn:  "1h",
			RefreshExpiresIn: "24h",
		},
	}
}
