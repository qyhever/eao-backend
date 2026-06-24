package model

type Disc struct {
	ImgURL  string `json:"imgURL"`
	Title   string `json:"title"`
	PlayNum string `json:"playNum"`
}

type DiscListQuery struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}

type DiscListResponse struct {
	List     []Disc `json:"list"`
	Total    int    `json:"total"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}
