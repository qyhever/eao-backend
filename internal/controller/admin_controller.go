package controller

import (
	"errors"
	"strconv"

	"eao/internal/model"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminController struct {
	adminService *service.AdminService
}

func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

func (ac *AdminController) GetAdmin(c *gin.Context) {
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

func (ac *AdminController) UpdateAdmin(c *gin.Context) {
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

func (ac *AdminController) BatchDeleteAdmins(c *gin.Context) {
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

func (ac *AdminController) ToggleAdminStatus(c *gin.Context) {
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

func (ac *AdminController) writeAdminError(c *gin.Context, logMsg string, err error) {
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
