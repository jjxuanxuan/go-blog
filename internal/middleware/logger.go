package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// LoggerMiddleware 打印简易访问日志：方法、路径、状态码、耗时。
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 请求前
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Printf("[START] %s %s\n", method, path)

		c.Next() // 继续处理请求

		// 请求后
		status := c.Writer.Status()
		latency := time.Since(start)
		fmt.Printf("[END] %s %s -> %d (%v)\n", method, path, status, latency)
	}
}
