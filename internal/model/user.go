package model

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=64"`
}
