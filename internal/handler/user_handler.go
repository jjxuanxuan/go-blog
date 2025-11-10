package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go-blog/internal/dto"
	"go-blog/internal/middleware"
	"go-blog/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// UserLogin 登入接口
func (h *UserHandler) UserLogin(c *gin.Context) {
	var req dto.LoginReq
	//绑定参数
    if err := c.ShouldBindJSON(&req); err != nil {
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
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
			"detail":  err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	token, err := middleware.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "签发 Token 失败"})
		return
	}
	//成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})

}

// UserRegister 注册接口
func (h *UserHandler) UserRegister(c *gin.Context) {
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	var count int64
	if err := h.DB.Model(&model.User{}).Where("username= ? or email= ?", req.Username, req.Email).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "数据库查询失败",
			"detail":  err.Error(),
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名或邮箱已存在",
		})
		return
	}

	// bcrypt 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	//用户模型
	u := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
	}

	//创建
	if err := h.DB.Create(&u).Error; err != nil {
		if isDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": "用户名或邮箱已存在",
			})
			return
		}
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
func (h *UserHandler) MeHandler(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
		return
	}
	var u model.User
	if err := h.DB.First(&u, v.(uint)).Error; err != nil {
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

func isDuplicateKeyError(err error) bool {
	// 方案A：精准判断 go-sql-driver 的 MySQLError
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}
	// 方案B：保底匹配（不同中间层/方言时）
	return strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
		strings.Contains(err.Error(), "1062")
}
