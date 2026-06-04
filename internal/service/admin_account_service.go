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
	ErrAdminUpdateFieldsRequired  = errors.New("username、name 或 password 至少提供一个")
	ErrAdminNotFound              = errors.New("管理员不存在")
	ErrAdminIDsRequired           = errors.New("管理员 id 列表不能为空")
	ErrAdminIDInvalid             = errors.New("管理员 id 必须大于 0")
	ErrAdminStatusInvalid         = errors.New("管理员状态只能是 active 或 disabled")
)

type AdminAccountService struct {
	adminRepo repository.AdminAccountRepository
}

func NewAdminAccountService(adminRepo repository.AdminAccountRepository) *AdminAccountService {
	return &AdminAccountService{
		adminRepo: adminRepo,
	}
}

func (s *AdminAccountService) CreateAdmin(ctx context.Context, req model.CreateAdminRequest) (*model.AdminProfileResponse, error) {

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

func (s *AdminAccountService) GetAdmin(ctx context.Context, id int64) (*model.AdminProfileResponse, error) {
	if id <= 0 {
		return nil, ErrAdminIDInvalid
	}

	admin, err := s.adminRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return adminProfileResponse(admin), nil
}

func (s *AdminAccountService) UpdateAdmin(ctx context.Context, id int64, req model.UpdateAdminRequest) (*model.AdminProfileResponse, error) {
	if id <= 0 {
		return nil, ErrAdminIDInvalid
	}
	if req.Username == nil && req.Name == nil && req.Password == nil {
		return nil, ErrAdminUpdateFieldsRequired
	}

	admin, err := s.adminRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return nil, ErrAdminUsernameRequired
		}
		admin.Username = username
	}
	if req.Name != nil {
		admin.Name = strings.TrimSpace(*req.Name)
	}
	if strings.TrimSpace(admin.Name) == "" {
		admin.Name = admin.Username
	}
	if req.Password != nil {
		rawPassword := strings.TrimSpace(*req.Password)
		if rawPassword == "" {
			return nil, ErrAdminPasswordRequired
		}
		hash, err := password.Hash(rawPassword)
		if err != nil {
			return nil, err
		}
		admin.PasswordHash = hash
	}

	updated, err := s.adminRepo.Update(ctx, *admin)
	if err != nil {
		if errors.Is(err, repository.ErrAdminUsernameAlreadyExists) {
			return nil, ErrAdminUsernameAlreadyExists
		}
		return nil, err
	}
	if updated == nil {
		return nil, ErrAdminNotFound
	}

	return adminProfileResponse(updated), nil
}

func (s *AdminAccountService) BatchDeleteAdmins(ctx context.Context, req model.BatchDeleteAdminRequest) error {
	if len(req.IDs) == 0 {
		return ErrAdminIDsRequired
	}
	for _, id := range req.IDs {
		if id <= 0 {
			return ErrAdminIDInvalid
		}
	}
	return s.adminRepo.BatchSoftDelete(ctx, req.IDs)
}

func (s *AdminAccountService) ToggleAdminStatus(ctx context.Context, id int64, req model.ToggleAdminStatusRequest) (*model.AdminProfileResponse, error) {
	if id <= 0 {
		return nil, ErrAdminIDInvalid
	}

	status := strings.TrimSpace(req.Status)
	if status != "active" && status != "disabled" {
		return nil, ErrAdminStatusInvalid
	}

	admin, err := s.adminRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return adminProfileResponse(admin), nil
}

func adminProfileResponse(admin *model.Admin) *model.AdminProfileResponse {
	return &model.AdminProfileResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Name:     admin.Name,
		Status:   admin.Status,
	}
}
