package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-blog/internal/dto"
	"go-blog/internal/middleware"
	"go-blog/internal/repository"
	"go-blog/internal/service"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PostHandler 处理文章相关 HTTP 请求。
type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
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

	uid := middleware.UID(c)

	post, err := h.svc.CreatePost(c.Request.Context(), uid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文章失败",
			"detail":  err.Error(),
		})
		return
	}

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
	posts, err := h.svc.GetAllPosts(c.Request.Context())
	if err != nil {
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
// 保留你原来的函数名，避免改路由。
func (h *PostHandler) GetPostsById(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	id := uint(id64)

	post, err := h.svc.GetPostByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPostNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询失败",
				"detail":  err.Error(),
			})
		}
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
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	id := uint(id64)

	var req dto.UpdatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	uid := middleware.UID(c)

	post, err := h.svc.UpdatePost(c.Request.Context(), uid, id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPostNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该文章"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新失败",
				"detail":  err.Error(),
			})
		}
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
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	id := uint(id64)

	uid := middleware.UID(c)

	if err := h.svc.DeletePost(c.Request.Context(), uid, id); err != nil {
		switch {
		case errors.Is(err, service.ErrPostNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该文章"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除失败",
				"detail":  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ListPosts 列表查询：分页 + 分类筛选 + 标签筛选 + 状态 + 关键字 + 排序
func (h *PostHandler) ListPosts(c *gin.Context) {
	// 解析 query 参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var categoryID *uint
	if cidStr := c.Query("category_id"); cidStr != "" {
		if cid64, err := strconv.ParseUint(cidStr, 10, 64); err == nil {
			cid := uint(cid64)
			categoryID = &cid
		}
	}

	// tag_ids=1,2,3
	var tagIDs []uint
	if tidStr := c.Query("tag_ids"); tidStr != "" {
		for _, s := range strings.Split(tidStr, ",") {
			if s == "" {
				continue
			}
			if id64, err := strconv.ParseUint(s, 10, 64); err == nil {
				tagIDs = append(tagIDs, uint(id64))
			}
		}
	}

	var status *string
	if st := c.Query("status"); st != "" {
		status = &st
	}

	keyword := c.Query("keyword")
	order := c.DefaultQuery("order", "latest")

	filter := repository.PostFilter{
		CategoryID: categoryID,
		TagIDs:     tagIDs,
		Status:     status,
		Keyword:    keyword,
		Order:      order,
		Page:       page,
		PageSize:   pageSize,
	}

	posts, total, err := h.svc.ListPosts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询文章失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":     posts,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
