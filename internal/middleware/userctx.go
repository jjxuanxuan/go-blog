package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const CtxUIDKey = "uid"

// RequireUser确保在上下文中存在一个有效的uid（由AuthMiddleware设置）。
// 如果缺少或无效，则使用401终止。
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

// UID从上下文返回类型化的UID
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
