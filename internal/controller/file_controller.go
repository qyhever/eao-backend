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

// List godoc
// @Summary 获取文件列表
// @Description 透传第三方文件服务的文件列表 JSON。
// @Tags file
// @Accept json
// @Produce json
// @Success 200 {object} SwaggerFileProxyResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/file/list [get]
func (fc *FileController) List(c *gin.Context) {
	body, err := fc.fileService.List(c.Request.Context())
	fc.writeProxyResponse(c, body, err)
}

// ListByDir godoc
// @Summary 按目录获取文件列表
// @Description 根据目录名称透传第三方文件服务的文件列表 JSON。
// @Tags file
// @Accept json
// @Produce json
// @Param dirName query string true "目录名称"
// @Success 200 {object} SwaggerFileProxyResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/file/listByDir [get]
func (fc *FileController) ListByDir(c *gin.Context) {
	body, err := fc.fileService.ListByDir(c.Request.Context(), c.Query("dirName"))
	fc.writeProxyResponse(c, body, err)
}

// Upload godoc
// @Summary 上传文件
// @Description 上传文件到指定目录，并透传第三方文件服务返回的 JSON。
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param dirName formData string true "目录名称"
// @Param file formData file true "文件"
// @Success 200 {object} SwaggerFileProxyResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/file/upload [post]
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
