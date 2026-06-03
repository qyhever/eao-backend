package service

import (
	"context"
	"errors"
	"testing"

	"eao/internal/model"
	"eao/internal/pkg/password"
	"eao/internal/repository"
	"eao/internal/repository/persistence"
)

func TestAdminServiceGetAdminSkipsSoftDeleted(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminService(repo)

	admin, err := svc.GetAdmin(ctx, 1)
	if err != nil {
		t.Fatalf("get admin failed: %v", err)
	}
	if admin.Username != "admin" {
		t.Fatalf("unexpected username: %s", admin.Username)
	}

	if err := svc.BatchDeleteAdmins(ctx, model.BatchDeleteAdminRequest{IDs: []int64{1}}); err != nil {
		t.Fatalf("delete admin failed: %v", err)
	}
	if _, err := svc.GetAdmin(ctx, 1); !errors.Is(err, ErrAdminNotFound) {
		t.Fatalf("expected ErrAdminNotFound, got %v", err)
	}
}

func TestAdminServiceUpdateAdmin(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminService(repo)

	username := "root"
	name := "超级管理员"
	rawPassword := "new-password"
	updated, err := svc.UpdateAdmin(ctx, 1, model.UpdateAdminRequest{
		Username: &username,
		Name:     &name,
		Password: &rawPassword,
	})
	if err != nil {
		t.Fatalf("update admin failed: %v", err)
	}
	if updated.Username != username || updated.Name != name {
		t.Fatalf("unexpected admin profile: %+v", updated)
	}

	admin, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find admin failed: %v", err)
	}
	if err := password.Compare(admin.PasswordHash, rawPassword); err != nil {
		t.Fatalf("password was not updated: %v", err)
	}
}

func TestAdminServiceUpdateAdminDuplicateUsername(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin 1 failed: %v", err)
	}
	if err := repo.Upsert(ctx, testAdmin(2, "root")); err != nil {
		t.Fatalf("seed admin 2 failed: %v", err)
	}
	svc := NewAdminService(repo)

	username := "root"
	_, err := svc.UpdateAdmin(ctx, 1, model.UpdateAdminRequest{Username: &username})
	if !errors.Is(err, ErrAdminUsernameAlreadyExists) {
		t.Fatalf("expected ErrAdminUsernameAlreadyExists, got %v", err)
	}

	admin, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("find admin failed: %v", err)
	}
	if admin.Username != "admin" {
		t.Fatalf("username should not change on duplicate, got %s", admin.Username)
	}
}

func TestAdminServiceBatchDeleteIsIdempotent(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminService(repo)

	req := model.BatchDeleteAdminRequest{IDs: []int64{1, 2}}
	if err := svc.BatchDeleteAdmins(ctx, req); err != nil {
		t.Fatalf("first delete failed: %v", err)
	}
	if err := svc.BatchDeleteAdmins(ctx, req); err != nil {
		t.Fatalf("second delete failed: %v", err)
	}
	if _, err := repo.FindByID(ctx, 1); err != nil {
		t.Fatalf("find admin failed: %v", err)
	} else if _, err := svc.GetAdmin(ctx, 1); !errors.Is(err, ErrAdminNotFound) {
		t.Fatalf("expected deleted admin not found, got %v", err)
	}
}

func TestAdminServiceToggleAdminStatus(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	svc := NewAdminService(repo)

	updated, err := svc.ToggleAdminStatus(ctx, 1, model.ToggleAdminStatusRequest{Status: "disabled"})
	if err != nil {
		t.Fatalf("disable admin failed: %v", err)
	}
	if updated.Status != "disabled" {
		t.Fatalf("unexpected status: %s", updated.Status)
	}

	if _, err := svc.ToggleAdminStatus(ctx, 1, model.ToggleAdminStatusRequest{Status: "locked"}); !errors.Is(err, ErrAdminStatusInvalid) {
		t.Fatalf("expected ErrAdminStatusInvalid, got %v", err)
	}
}

func TestAdminRepositoryFindByUsernameSkipsSoftDeleted(t *testing.T) {
	ctx := context.Background()
	repo := persistence.NewAdminRepository(nil)
	if err := repo.Upsert(ctx, testAdmin(1, "admin")); err != nil {
		t.Fatalf("seed admin failed: %v", err)
	}
	if err := repo.BatchSoftDelete(ctx, []int64{1}); err != nil {
		t.Fatalf("delete admin failed: %v", err)
	}

	admin, err := repo.FindByUsername(ctx, "admin")
	if err != nil {
		t.Fatalf("find by username failed: %v", err)
	}
	if admin != nil {
		t.Fatalf("expected nil for soft deleted admin, got %+v", admin)
	}

	if _, err := repo.Create(ctx, testAdmin(0, "admin")); !errors.Is(err, repository.ErrAdminUsernameAlreadyExists) {
		t.Fatalf("expected duplicate username for soft deleted admin, got %v", err)
	}
}

func testAdmin(id int64, username string) model.Admin {
	hash, err := password.Hash("password")
	if err != nil {
		panic(err)
	}
	return model.Admin{
		ID:           id,
		Username:     username,
		PasswordHash: hash,
		Name:         username,
		Status:       "active",
	}
}
