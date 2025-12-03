package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-blog/internal/dto"
	"go-blog/internal/service"
	"go-blog/internal/util"
)

// AdminHandler 处理后台管理相关 HTTP 请求。
type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// Dashboard 管理端仪表盘数据。
func (h *AdminHandler) Dashboard(c *gin.Context) {
	topStr := c.DefaultQuery("top", "5")
	top, err := strconv.Atoi(topStr)
	if err != nil || top <= 0 {
		top = 5
	}
	if top > 20 {
		top = 20
	}

	dashboard, err := h.svc.Dashboard(c.Request.Context(), top)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取仪表盘数据失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dashboard,
	})
}

// ListUsers 管理端分页查询用户。
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, pageSize := util.ParsePage(c)
	query := dto.AdminUserQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Role:     c.Query("role"),
	}

	users, total, err := h.svc.ListUsers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": util.PageResult{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     users,
		},
	})
}

// ListPosts 管理端分页查询文章。
func (h *AdminHandler) ListPosts(c *gin.Context) {
	page, pageSize := util.ParsePage(c)
	var status *string
	if st := c.Query("status"); st != "" {
		status = &st
	}
	query := dto.AdminPostQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
		Status:   status,
	}

	posts, total, err := h.svc.ListPosts(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询文章失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": util.PageResult{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     posts,
		},
	})
}

// ListComments 管理端分页查询评论。
func (h *AdminHandler) ListComments(c *gin.Context) {
	page, pageSize := util.ParsePage(c)
	query := dto.AdminCommentQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  c.Query("keyword"),
	}

	if uidStr := c.Query("user_id"); uidStr != "" {
		if uid64, err := strconv.ParseUint(uidStr, 10, 64); err == nil && uid64 > 0 {
			uid := uint(uid64)
			query.UserID = &uid
		}
	}

	if pidStr := c.Query("post_id"); pidStr != "" {
		if pid64, err := strconv.ParseUint(pidStr, 10, 64); err == nil && pid64 > 0 {
			pid := uint(pid64)
			query.PostID = &pid
		}
	}

	comments, total, err := h.svc.ListComments(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论失败",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": util.PageResult{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     comments,
		},
	})
}
