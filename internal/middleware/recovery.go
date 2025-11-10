package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "runtime/debug"
)

// Recovery 中间件在发生 panic 时返回 JSON 并防止进程崩溃。
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // 记录堆栈跟踪到 stderr；在实际应用中集成到日志系统
                _ = debug.Stack()
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code":    500,
                    "message": "服务器内部错误",
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}

