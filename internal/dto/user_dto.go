package dto

// CreateUserReq 用于注册用户的请求体
type CreateUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// LoginReq 用于用户登录
type LoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// UpdateUserReq 用于更新资料
type UpdateUserReq struct {
	Email    *string `json:"email"    binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6,max=64"`
}
