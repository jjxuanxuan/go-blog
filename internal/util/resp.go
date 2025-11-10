// Package util 提供通用的工具函数：统一响应、哈希、JWT 等。
package util

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Response 统一的响应结构，所有业务接口尽量返回该格式。
// code: 0 表示成功，非 0 表示业务错误；HTTP 状态码体现协议级错误。
type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Detail    string      `json:"detail,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
}

// OK 快捷成功响应，message 固定为 "ok"。
func OK(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Success 自定义成功 message 的快捷方法。
func Success(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 0, Message: message, Data: data})
}

// Fail 失败响应，httpStatus 为 HTTP 状态码，code 为业务码，detail 可选。
func Fail(c *gin.Context, httpStatus, code int, message string, detail string) {
    c.JSON(httpStatus, Response{Code: code, Message: message, Detail: detail})
}
