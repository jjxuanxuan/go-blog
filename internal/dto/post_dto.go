package dto

// CreatePostReq 用于创建文章请求体
type CreatePostReq struct {
	Title   string `json:"title"   binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required"`
	UserID  uint   `json:"user_id" binding:"required"`
}

// UpdatePostReq 用于更新文章请求体
type UpdatePostReq struct {
	Title   *string `json:"title"   binding:"omitempty,min=1,max=200"`
	Content *string `json:"content" binding:"omitempty"`
}
