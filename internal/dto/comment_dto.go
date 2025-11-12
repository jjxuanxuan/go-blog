package dto

// CreateCommentReq 创建评论请求
type CreateCommentReq struct{
	Content string `json:"content" binding:"required,min=1,max=1000"`
	PostId uint `json:"post_id" binding:"required"`
}

// ReplyCommentReq 回复评论请求
type ReplyCommentReq struct{
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// UserBrief 用户简要信息
type UserBrief struct{
	Id uint `json:"id"`
	Username string `json:"username"`
}

// CommentResp 评论响应
type CommentResp struct{
	Id uint `json:"id"`
	Content string `json:"content"`
	User UserBrief `json:"user"`
	ParentId *uint `json:"parent_id,omitempty"`
	PostId uint `json:"post_id"`
	Replies []CommentResp `json:"replies,omitempty"`
}