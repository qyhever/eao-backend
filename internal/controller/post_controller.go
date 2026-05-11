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
