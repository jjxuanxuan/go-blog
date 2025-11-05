package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-blog/internal/dto"
	"go-blog/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// LoginHandler 登入接口
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginReq
	//绑定参数
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	var u model.User
	if err := h.DB.Where("username=?", req.Username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
			"detail":  err.Error(),
		})
		return
	}

	if u.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	//成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"username": u.Username,
			"token":    "dummy-token-" + u.Username,
		},
	})

}

// PingHandler PingHandle 检查接口健康
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// RegisterHandler 注册接口
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	var existing model.User

	//用户名或者邮箱已存在返回401
	//没找到用户或者邮箱 注册
	//出现其他错误返回500
	if err := h.DB.Where("username= ? or email= ?", req.Username, req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或者邮箱已存在",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "数据库查询失败",
			"detail":  err.Error(),
		})
	}

	u := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		},
	})
}
