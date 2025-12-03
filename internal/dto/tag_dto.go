package dto

// TagResp 标签响应体。
type TagResp struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Weight int    `json:"weight"`
}

// CreateTagReq 创建标签请求体。
type CreateTagReq struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	Slug string `json:"slug" binding:"required,min=1,max=100"`
}
