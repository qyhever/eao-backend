package model

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type PostListQuery struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
}

type PostListResponse struct {
	List     []Post `json:"list"`
	Total    int    `json:"total"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

type PostDetailResponse struct {
	Post *Post `json:"post"`
}

type PostMutationResponse struct {
	ID string `json:"id"`
}
