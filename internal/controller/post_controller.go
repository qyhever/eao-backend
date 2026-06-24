package controller

import (
	"eao/internal/model"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostController struct {
	postService *service.PostService
}

func NewPostController(postService *service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// GetPostList godoc
// @Summary 获取文章列表
// @Description 分页获取文章列表，可按关键字过滤。
// @Tags post
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键字"
// @Success 200 {object} SwaggerPostListResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/post [get]
func (pc *PostController) GetPostList(c *gin.Context) {
	var req model.PostListQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	list, total, err := pc.postService.GetPostList(&req)
	if err != nil {
		zap.L().Error("get post list failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, &model.PostListResponse{
		List:     list,
		Total:    total,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
}

// GetPostByID godoc
// @Summary 获取文章详情
// @Description 根据文章 ID 获取文章详情。
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "文章 ID"
// @Success 200 {object} SwaggerPostDetailResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/post/{id} [get]
func (pc *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ResponseFailedWithMsg(c, CodeInvalidParam, "id 不能为空")
		return
	}

	post, err := pc.postService.GetPostByID(id)
	if err != nil {
		zap.L().Error("get post by id failed", zap.String("id", id), zap.Error(err))
		ResponseFailedWithMsg(c, CodeResourceNotExist, err.Error())
		return
	}

	ResponseSuccess(c, &model.PostDetailResponse{Post: post})
}

// CreatePost godoc
// @Summary 创建文章
// @Description 创建一篇文章并返回文章 ID。
// @Tags post
// @Accept json
// @Produce json
// @Param request body model.CreatePostRequest true "文章内容"
// @Success 200 {object} SwaggerPostMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/post [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	id, err := pc.postService.CreatePost(&req)
	if err != nil {
		zap.L().Error("create post failed", zap.Error(err))
		ResponseFailedWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	ResponseSuccess(c, &model.PostMutationResponse{ID: id})
}

// UpdatePost godoc
// @Summary 更新文章
// @Description 根据文章 ID 更新标题或内容，title 和 content 至少提供一个。
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "文章 ID"
// @Param request body model.UpdatePostRequest true "文章更新内容"
// @Success 200 {object} SwaggerPostMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/post/{id} [put]
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ResponseFailedWithMsg(c, CodeInvalidParam, "id 不能为空")
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "请求参数错误: "+err.Error())
		return
	}

	if req.Title == nil && req.Content == nil {
		ResponseFailedWithMsg(c, CodeInvalidParam, "title 或 content 至少提供一个")
		return
	}

	err := pc.postService.UpdatePost(id, &req)
	if err != nil {
		zap.L().Error("update post failed", zap.String("id", id), zap.Error(err))
		ResponseFailedWithMsg(c, CodeResourceNotExist, err.Error())
		return
	}

	ResponseSuccess(c, &model.PostMutationResponse{ID: id})
}

// DeletePost godoc
// @Summary 删除文章
// @Description 根据文章 ID 删除文章。
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "文章 ID"
// @Success 200 {object} SwaggerPostMutationResponse
// @Failure 200 {object} SwaggerErrorResponse
// @Router /common/post/{id} [delete]
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ResponseFailedWithMsg(c, CodeInvalidParam, "id 不能为空")
		return
	}

	err := pc.postService.DeletePost(id)
	if err != nil {
		zap.L().Error("delete post failed", zap.String("id", id), zap.Error(err))
		ResponseFailedWithMsg(c, CodeResourceNotExist, err.Error())
		return
	}

	ResponseSuccess(c, &model.PostMutationResponse{ID: id})
}
