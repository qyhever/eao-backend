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

func (ac *AdminAuthController) writeAdminAuthError(c *gin.Context, logMsg string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidAdminCredentials):
		ResponseFailedWithMsg(c, CodeInvalidPassword, err.Error())
	case errors.Is(err, service.ErrAdminUsernameRequired),
		errors.Is(err, service.ErrAdminPasswordRequired):
		ResponseFailedWithMsg(c, CodeInvalidParam, err.Error())
	default:
		zap.L().Error(logMsg, zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
	}
}
