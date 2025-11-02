package handler

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/model"
	"net/http"
)

// 模拟用户数据
var users = map[string]string{
	"admin":     "123456",
	"xiaojiang": "030722",
}

// LoginHandle 登入接口
func LoginHandler(c *gin.Context) {
	var req model.LoginRequest

	//绑定参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "参数错误",
			"defail":  err.Error(),
		})
		return
	}

	//用户名是否存在；
	//密码是否正确。
	password, ok := users[req.Username]
	if !ok || password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    2001,
			"message": "用户名或密码错误",
		})
		return
	}

	//成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登入成功",
		"data": gin.H{
			"username": req.Username,
			"token":    "dummy-token-" + req.Username,
		},
	})

}

// PingHandler PingHandle 检查接口健康
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
