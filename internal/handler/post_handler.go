package handler

import (
    "errors"
    "github.com/gin-gonic/gin"
    "go-blog/internal/dto"
    "go-blog/internal/middleware"
    "go-blog/internal/model"
    "gorm.io/gorm"
    "net/http"
    "time"
)

type PostHandler struct {
	DB *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{DB: db}
}

// CreatePost 创建文章：从上下文获取 uid，避免客户端伪造 user_id。
func (h *PostHandler) CreatePost(c *gin.Context) {
    var req dto.CreatePostReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "创建文章失败",
            "detail":  err.Error(),
        })
        return
    }

    // 从中间件中获取 uid（RequireUser 已保证存在）
    uid := middleware.UID(c)

    post := model.Post{
        Title:   req.Title,
        Content: req.Content,
        UserID:  uid,
    }

    if err := h.DB.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "创建文章失败",
            "detail":  err.Error(),
        })
        return
    }

	//创建成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建文章成功",
		"data": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"created_at": post.CreatedAt.Format(time.RFC3339),
		},
	})

}

// GetAllPosts 获取所有文章列表（预加载作者信息）。
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	var posts []model.Post

	//查询所有文章和它的作者
	if err := h.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询成功",
		"data":    posts,
	})

}

// GetPostsById 根据ID查询单篇文章详情，预加载作者信息。
func (h *PostHandler) GetPostsById(c *gin.Context) {
	id := c.Param("id")

	var post model.Post
	if err := h.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询成功",
		"data":    post,
	})
}

// UpdatePost 更新文章内容：仅作者本人可更新，空字段不覆盖。
func (h *PostHandler) UpdatePost(c *gin.Context) {

    id := c.Param("id")
    var req dto.UpdatePostReq
    if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

    // 鉴权用户（RequireUser 已保证存在）
    uid := middleware.UID(c)

    var post model.Post
    if err := h.DB.First(&post, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "detail": err.Error()})
        return
    }

    if post.UserID != uid {
        c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该文章"})
        return
    }

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
    // 让 GORM 自动维护 UpdatedAt，不要覆盖 CreatedAt

    if err := h.DB.Save(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "更新失败",
            "detail":  err.Error(),
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    post,
	})

}

// DeletePost 删除文章：仅作者本人可删除。
func (h *PostHandler) DeletePost(c *gin.Context) {

    id := c.Param("id")

    // 鉴权用户（RequireUser 已保证存在）
    uid := middleware.UID(c)

    var post model.Post
    if err := h.DB.First(&post, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "detail": err.Error()})
        return
    }

    if post.UserID != uid {
        c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该文章"})
        return
    }

    if err := h.DB.Delete(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "删除失败",
            "detail":  err.Error(),
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
