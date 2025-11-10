package handler

import (
    "errors"
    "github.com/gin-gonic/gin"
    "go-blog/internal/dto"
    "go-blog/internal/model"
    "go-blog/internal/util"
    "gorm.io/gorm"
    "net/http"
)

type AuthHandler struct{
    DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler { return &AuthHandler{DB: db} }

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.CreateUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":400, 
            "message":"参数错误", 
            "detail": err.Error(),
        })
        return
    }
    
    // 检查用户名或邮箱是否已存在
    var count int64
    if err := h.DB.Model(&model.User{}).Where("username=? OR email=?", req.Username, req.Email).Count(&count).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"数据库查询失败", 
            "detail": err.Error(),
        })
        return
    }
    if count > 0 {
        c.JSON(http.StatusConflict, gin.H{
            "code":409, 
            "message":"用户名或邮箱已存在",
        })
        return
    }

    hashed, err := util.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"密码加密失败",
        })
        return
    }

    u := model.User{
        Username: req.Username, 
        Email: req.Email, 
        Password: hashed, 
        Role: "user",
    }
    if err := h.DB.Create(&u).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"创建用户失败", 
            "detail": err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "code":0, 
        "message":"注册成功", 
        "data": gin.H{
            "id": u.ID, 
            "username": u.Username, 
            "email": u.Email,
        },
    })
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":400, 
            "message":"参数错误", 
            "detail": err.Error(),
        })
        return
    }
    var u model.User
    if err := h.DB.Where("username=?", req.Username).First(&u).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code":401, 
                "message":"用户名或密码错误",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"查询用户失败", 
            "detail": err.Error(),
        })
        return
    }
    if !util.CheckPassword(u.Password, req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code":401, 
            "message":"用户名或密码错误",
        })
        return
    }
    at, err := util.GenerateAccessToken(u.ID, u.Role)
    if err != nil { 
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"签发 access_token 失败",
        }); return 
    }
    rt, err := util.GenerateRefreshToken(u.ID, u.Role)
    if err != nil { 
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"签发 refresh_token 失败",
        }); return 
    }
    c.JSON(http.StatusOK, gin.H{
        "code":0, 
        "message":"登录成功", 
        "access_token": at, 
        "refresh_token": rt,
    })
}

type refreshReq struct { RefreshToken string `json:"refresh_token" binding:"required"` }

// POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
    var req refreshReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":400, 
            "message":"参数错误", 
            "detail": err.Error(),
        })
        return
    }
    claims, err := util.ParseToken(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code":401, 
            "message":"无效Token",
        })
        return
    }
    // subject 是用户ID
    // 将 subject 转换为 uint
    var uid uint
    if claims.Subject != "" {
        // 安全解析；忽略错误 -> uid=0 表示无效
        var n uint64
        for i := 0; i < len(claims.Subject); i++ { // 轻量级 atoi 实现
            ch := claims.Subject[i]
            if ch < '0' || ch > '9' { n = 0; break }
            n = n*10 + uint64(ch-'0')
        }
        uid = uint(n)
    }
    if uid == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code":401, 
            "message":"无效Token",
        })
        return
    }
    at, err := util.GenerateAccessToken(uid, claims.Role)
    if err != nil { 
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":500, 
            "message":"签发 access_token 失败",
        }); return 
    }
    c.JSON(http.StatusOK, gin.H{
        "code":0, 
        "message":"ok", 
        "access_token": at,
    })
}

