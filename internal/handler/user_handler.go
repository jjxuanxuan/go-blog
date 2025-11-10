// Package handler/user 提供与用户相关的业务接口。
package handler

import (
    "github.com/gin-gonic/gin"
    "go-blog/internal/middleware"
    "go-blog/internal/model"
    "gorm.io/gorm"
    "net/http"
)

// UserHandler 处理用户相关的 HTTP 请求。
type UserHandler struct {
    DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler { return &UserHandler{DB: db} }
// MeHandler 返回当前登录用户的基本信息。
func (h *UserHandler) MeHandler(c *gin.Context) {
    uid := middleware.UID(c)
    var u model.User
    if err := h.DB.First(&u, uid).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "ok",
        "data": gin.H{
            "id":       u.ID,
            "username": u.Username,
            "email":    u.Email,
        },
    })
}
