package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CtxUIDKey 在上下文中保存用户ID的键名。
const CtxUIDKey = "uid"

// RequireUser 确保上下文中存在有效 uid（由 AuthMiddleware 设置），缺失则返回 401。
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
			c.Abort()
			return
		}
		uid, ok := v.(uint)
		if !ok || uid == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
			c.Abort()
			return
		}
		c.Set(CtxUIDKey, uid)
		c.Next()
	}
}

// UID 从上下文返回类型化的 UID。
func UID(c *gin.Context) uint {
	if v, ok := c.Get(CtxUIDKey); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	if v, ok := c.Get("user_id"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}
