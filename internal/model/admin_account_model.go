package model

import "time"

type Admin struct {
	ID           int64
	Username     string
	PasswordHash string
	Name         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
}

type AdminProfileResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type UpdateAdminRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Name     *string `json:"name"`
}

type BatchDeleteAdminRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

type ToggleAdminStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminMutationResponse struct {
	ID int64 `json:"id"`
}

type AdminBatchMutationResponse struct {
	IDs []int64 `json:"ids"`
}
