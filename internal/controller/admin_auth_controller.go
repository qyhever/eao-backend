package controller

import (
	"errors"

	"eao/internal/model"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminAuthController struct {
	authService *service.AdminAuthService
}

func NewAdminAuthController(
	authService *service.AdminAuthService,
) *AdminAuthController {
	return &AdminAuthController{
		authService: authService,
	}
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 使用管理员用户名和密码换取 accessToken 与 refreshToken。
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param request body model.AdminLoginRequest true "登录参数"
// @Success 200 {object} SwaggerAdminLoginResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/auth/login [post]
func (ac *AdminAuthController) AdminLogin(c *gin.Context) {
	var param model.AdminLoginRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	result, err := ac.authService.AdminLogin(c.Request.Context(), param)
	if err != nil {
		ac.writeAdminAuthError(c, "admin login failed", err)
		return
	}

	ResponseSuccess(c, result)
}

// AdminRefreshToken godoc
// @Summary 刷新管理员 Token
// @Description 使用 refreshToken 换取新的 accessToken 与 refreshToken。
// @Tags admin-auth
// @Accept json
// @Produce json
// @Param request body model.AdminRefreshTokenRequest true "刷新参数"
// @Success 200 {object} SwaggerAdminLoginResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/auth/refresh [post]
func (ac *AdminAuthController) AdminRefreshToken(c *gin.Context) {
	var param model.AdminRefreshTokenRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, err.Error())
		return
	}

	result, err := ac.authService.AdminRefreshToken(c.Request.Context(), param)
	if err != nil {
		ac.writeAdminAuthError(c, "admin refresh token failed", err)
		return
	}

	ResponseSuccess(c, result)
}

func (ac *AdminAuthController) writeAdminAuthError(c *gin.Context, logMsg string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidAdminCredentials):
		ResponseFailed(c, CodeInvalidPassword)
	case errors.Is(err, service.ErrAdminUsernameRequired),
		errors.Is(err, service.ErrAdminPasswordRequired):
		ResponseFailedWithMsg(c, CodeInvalidParam, err.Error())
	case errors.Is(err, service.ErrInvalidRefreshToken):
		ResponseFailed(c, CodeInvalidToken)
	case errors.Is(err, service.ErrUserNotFound):
		ResponseFailed(c, CodeUserNotExist)
	default:
		zap.L().Error(logMsg, zap.Error(err))
		ResponseFailed(c, CodeServerBusy)
	}
}
