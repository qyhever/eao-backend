package controller

import "eao/internal/model"

type SwaggerErrorResponse struct {
	Code    MyCode      `json:"code" example:"1001"`
	Message string      `json:"message" example:"请求参数错误"`
	Data    interface{} `json:"data"`
}

type SwaggerHelloInfoResponse struct {
	Code    MyCode                     `json:"code" example:"1000"`
	Message string                     `json:"message" example:"success"`
	Data    model.GetHelloInfoResponse `json:"data"`
}

type SwaggerDiscListResponse struct {
	Code    MyCode                 `json:"code" example:"1000"`
	Message string                 `json:"message" example:"success"`
	Data    model.DiscListResponse `json:"data"`
}

type SwaggerPostListResponse struct {
	Code    MyCode                 `json:"code" example:"1000"`
	Message string                 `json:"message" example:"success"`
	Data    model.PostListResponse `json:"data"`
}

type SwaggerPostDetailResponse struct {
	Code    MyCode                   `json:"code" example:"1000"`
	Message string                   `json:"message" example:"success"`
	Data    model.PostDetailResponse `json:"data"`
}

type SwaggerPostMutationResponse struct {
	Code    MyCode                     `json:"code" example:"1000"`
	Message string                     `json:"message" example:"success"`
	Data    model.PostMutationResponse `json:"data"`
}

type SwaggerFileProxyResponse struct {
	Code    int                    `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type SwaggerVideoListResponse struct {
	Code    MyCode              `json:"code" example:"1000"`
	Message string              `json:"message" example:"success"`
	Data    []model.VideoConfig `json:"data"`
}

type SwaggerAdminLoginResponse struct {
	Code    MyCode                   `json:"code" example:"1000"`
	Message string                   `json:"message" example:"success"`
	Data    model.AdminLoginResponse `json:"data"`
}

type SwaggerAdminProfileResponse struct {
	Code    MyCode                     `json:"code" example:"1000"`
	Message string                     `json:"message" example:"success"`
	Data    model.AdminProfileResponse `json:"data"`
}

type SwaggerAdminMutationResponse struct {
	Code    MyCode                      `json:"code" example:"1000"`
	Message string                      `json:"message" example:"success"`
	Data    model.AdminMutationResponse `json:"data"`
}

type SwaggerAdminBatchMutationResponse struct {
	Code    MyCode                           `json:"code" example:"1000"`
	Message string                           `json:"message" example:"success"`
	Data    model.AdminBatchMutationResponse `json:"data"`
}
