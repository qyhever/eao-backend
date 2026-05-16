package controller

import (
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoController struct {
	videoService *service.VideoService
}

func NewVideoController(videoService *service.VideoService) *VideoController {
	return &VideoController{videoService: videoService}
}

func (vc *VideoController) GetVideoList(c *gin.Context) {
	list, err := vc.videoService.GetVideoList()
	if err != nil {
		zap.L().Error("get video list failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, list)
}
