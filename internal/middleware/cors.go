package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// CORS 设置跨域响应头，放行常见方法与头；对预检请求直接返回 204。
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.Writer.Header()
        h.Set("Access-Control-Allow-Origin", c.GetHeader("Origin"))
        h.Set("Vary", "Origin")
        h.Set("Access-Control-Allow-Credentials", "true")
        h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
        h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}
