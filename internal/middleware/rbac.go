package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// RequireRole 确保当前用户拥有允许的角色之一。
func RequireRole(roles ...string) gin.HandlerFunc {
    allowed := map[string]struct{}{}
    for _, r := range roles {
        allowed[r] = struct{}{}
    }
    return func(c *gin.Context) {
        v, ok := c.Get("role")
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
            c.Abort()
            return
        }
        role, ok := v.(string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
            c.Abort()
            return
        }
        if _, ok := allowed[role]; !ok {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
            c.Abort()
            return
        }
        c.Next()
    }
}

