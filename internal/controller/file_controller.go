package controller

import (
	"errors"
	"net/http"

	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileController struct {
	fileService *service.FileService
}

func NewFileController(fileService *service.FileService) *FileController {
	return &FileController{fileService: fileService}
}

func (fc *FileController) List(c *gin.Context) {
	body, err := fc.fileService.List(c.Request.Context())
	fc.writeProxyResponse(c, body, err)
}

func (fc *FileController) ListByDir(c *gin.Context) {
	body, err := fc.fileService.ListByDir(c.Request.Context(), c.Query("dirName"))
	fc.writeProxyResponse(c, body, err)
}

func (fc *FileController) Upload(c *gin.Context) {
	dirName := c.PostForm("dirName")
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "file 不能为空")
		return
	}
	defer file.Close()

	body, err := fc.fileService.Upload(c.Request.Context(), file, fileHeader, dirName)
	fc.writeProxyResponse(c, body, err)
}

func (fc *FileController) writeProxyResponse(c *gin.Context, body []byte, err error) {
	if err != nil {
		if errors.Is(err, service.ErrFileInvalidParam) {
			ResponseFailed(c, CodeInvalidParam)
			return
		}
		zap.L().Error("proxy file api failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}
