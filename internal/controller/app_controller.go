package controller

import (
	"eao/internal/model"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppController struct {
	appService *service.AppService
}

func NewAppController(appService *service.AppService) *AppController {
	return &AppController{
		appService: appService,
	}
}

// GetHelloInfo godoc
// @Summary 获取问候信息
// @Description 根据传入 name 返回问候信息。
// @Tags app
// @Accept json
// @Produce json
// @Param request body model.GetHelloInfoRequest true "请求参数"
// @Success 200 {object} SwaggerHelloInfoResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /app/getHelloInfo [post]
func (app *AppController) GetHelloInfo(c *gin.Context) {
	var req model.GetHelloInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	result, err := app.appService.GetHelloInfo(&req)
	if err != nil {
		zap.L().Error("get hello info failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, result)
}
