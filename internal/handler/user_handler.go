// Package handler/user 提供与用户相关的业务接口。
package handler

import (
	"errors"
	"go-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-blog/internal/middleware"
)

// UserHandler 处理用户相关的 HTTP 请求。
type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc: svc} }

// MeHandler 返回当前登录用户的基本信息。
func (h *UserHandler) MeHandler(c *gin.Context) {
	uid := middleware.UID(c)
	u, err := h.svc.GetMe(c.Request.Context(), uid)
	if err != nil {
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

// ListUserPosts 返回指定用户（仅限本人）的文章列表
func (h *UserHandler) ListUserPosts(c *gin.Context) {
	idStr := c.Param("id")
	uid64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || uid64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	requesterID := middleware.UID(c)
	targetUserID := uint(uid64)
	posts, err := h.svc.ListUserPosts(c.Request.Context(), requesterID, targetUserID)
	if err != nil {
		if errors.Is(err, service.ErrorForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权查看他人文章"})
			return
		}
		// 其他错误：系统错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户文章失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询成功",
		"data":    posts,
	})
}
