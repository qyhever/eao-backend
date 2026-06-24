package controller

import (
	"eao/internal/model"
	"eao/internal/pkg/pagination"
	"eao/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DiscController struct {
	discService *service.DiscService
}

func NewDiscController(discService *service.DiscService) *DiscController {
	return &DiscController{
		discService: discService,
	}
}

func (dc *DiscController) GetDiscList(c *gin.Context) {
	req := model.DiscListQuery{
		PageNum:  parseQueryInt(c, "pageNum"),
		PageSize: parseQueryInt(c, "pageSize"),
	}

	params := pagination.Normalize(req.PageNum, req.PageSize)
	req.PageNum = params.PageNum
	req.PageSize = params.PageSize

	list, total, err := dc.discService.GetDiscList(&req)
	if err != nil {
		zap.L().Error("get disc list failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, &model.DiscListResponse{
		List:     list,
		Total:    total,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
}

func parseQueryInt(c *gin.Context, key string) int {
	value := c.Query(key)
	if value == "" {
		return 0
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return num
}
