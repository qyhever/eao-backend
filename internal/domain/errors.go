package domain

import "errors"

var (
	ErrInvalidAdminCredentials = errors.New("账号或密码错误")
	ErrInvalidRefreshToken     = errors.New("refresh token 无效")
)
