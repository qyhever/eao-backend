package controller

import (
	"errors"
	"strconv"

	"eao/internal/model"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminAccountController struct {
	adminService *service.AdminAccountService
}

func NewAdminAccountController(adminService *service.AdminAccountService) *AdminAccountController {
	return &AdminAccountController{
		adminService: adminService,
	}
}

// GetAdmin godoc
// @Summary 获取管理员资料
// @Description 根据管理员 ID 获取管理员资料。Authorization 格式为 Bearer <accessToken>。
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员 ID"
// @Success 200 {object} SwaggerAdminProfileResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/users/{id} [get]
func (ac *AdminAccountController) GetAdmin(c *gin.Context) {
	id, ok := bindAdminID(c)
	if !ok {
		return
	}

	admin, err := ac.adminService.GetAdmin(c.Request.Context(), id)
	if err != nil {
		ac.writeAdminError(c, "get admin failed", err)
		return
	}

	ResponseSuccess(c, admin)
}

// UpdateAdmin godoc
// @Summary 更新管理员
// @Description 根据管理员 ID 更新管理员账号、密码或名称。Authorization 格式为 Bearer <accessToken>。
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员 ID"
// @Param request body model.UpdateAdminRequest true "更新参数"
// @Success 200 {object} SwaggerAdminMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/users/{id} [put]
func (ac *AdminAccountController) UpdateAdmin(c *gin.Context) {
	id, ok := bindAdminID(c)
	if !ok {
		return
	}

	var req model.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	admin, err := ac.adminService.UpdateAdmin(c.Request.Context(), id, req)
	if err != nil {
		ac.writeAdminError(c, "update admin failed", err)
		return
	}

	ResponseSuccess(c, &model.AdminMutationResponse{ID: admin.ID})
}

// BatchDeleteAdmins godoc
// @Summary 批量删除管理员
// @Description 根据管理员 ID 列表批量软删除管理员。Authorization 格式为 Bearer <accessToken>。
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.BatchDeleteAdminRequest true "批量删除参数"
// @Success 200 {object} SwaggerAdminBatchMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/users/batch [delete]
func (ac *AdminAccountController) BatchDeleteAdmins(c *gin.Context) {
	var req model.BatchDeleteAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	if err := ac.adminService.BatchDeleteAdmins(c.Request.Context(), req); err != nil {
		ac.writeAdminError(c, "batch delete admins failed", err)
		return
	}

	ResponseSuccess(c, &model.AdminBatchMutationResponse{IDs: req.IDs})
}

// ToggleAdminStatus godoc
// @Summary 切换管理员状态
// @Description 根据管理员 ID 切换管理员状态，status 支持 active 或 disabled。Authorization 格式为 Bearer <accessToken>。
// @Tags admin-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "管理员 ID"
// @Param request body model.ToggleAdminStatusRequest true "状态参数"
// @Success 200 {object} SwaggerAdminMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /admin/users/{id}/status [put]
func (ac *AdminAccountController) ToggleAdminStatus(c *gin.Context) {
	id, ok := bindAdminID(c)
	if !ok {
		return
	}

	var req model.ToggleAdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	admin, err := ac.adminService.ToggleAdminStatus(c.Request.Context(), id, req)
	if err != nil {
		ac.writeAdminError(c, "toggle admin status failed", err)
		return
	}

	ResponseSuccess(c, &model.AdminMutationResponse{ID: admin.ID})
}

func bindAdminID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ResponseFailedWithMsg(c, CodeInvalidParam, service.ErrAdminIDInvalid.Error())
		return 0, false
	}
	return id, true
}

func (ac *AdminAccountController) writeAdminError(c *gin.Context, logMsg string, err error) {
	switch {
	case errors.Is(err, service.ErrAdminIDInvalid),
		errors.Is(err, service.ErrAdminUsernameRequired),
		errors.Is(err, service.ErrAdminPasswordRequired),
		errors.Is(err, service.ErrAdminUpdateFieldsRequired),
		errors.Is(err, service.ErrAdminIDsRequired),
		errors.Is(err, service.ErrAdminStatusInvalid):
		ResponseFailedWithMsg(c, CodeInvalidParam, err.Error())
	case errors.Is(err, service.ErrAdminUsernameAlreadyExists):
		ResponseFailedWithMsg(c, CodeResourceExists, err.Error())
	case errors.Is(err, service.ErrAdminNotFound):
		ResponseFailedWithMsg(c, CodeResourceNotExist, err.Error())
	default:
		zap.L().Error(logMsg, zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
	}
}
