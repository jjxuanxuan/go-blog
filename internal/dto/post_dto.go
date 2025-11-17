package dto

// CreatePostReq 用于创建文章请求体
type CreatePostReq struct {
	Title      string `json:"title"   binding:"required,min=1,max=200"`
	Content    string `json:"content" binding:"required"`
	CategoryId uint   `json:"category_id" binding:"required"`
	TagIds     []uint `json:"tag_ids"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published"`
}

// UpdatePostReq 用于更新文章请求体
type UpdatePostReq struct {
	Title      *string `json:"title"   binding:"omitempty,min=1,max=200"`
	Content    *string `json:"content" binding:"omitempty"`
	CategoryID *uint   `json:"category_id" binding:"omitempty,gt=0"`                   // 分类可选更新
	Status     *string `json:"status"      binding:"omitempty,oneof=draft published"` // 状态：draft / published
	TagIDs     []uint  `json:"tag_ids"`
}
