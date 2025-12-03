package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-blog/internal/dto"
	"go-blog/internal/middleware"
	"go-blog/internal/service"
	"go-blog/internal/util"
)

// CommentHandler 处理评论相关 HTTP 接口。
type CommentHandler struct{ svc *service.CommentService }

func NewCommentHandler(svc *service.CommentService) *CommentHandler { return &CommentHandler{svc: svc} }

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req dto.CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}
	uid := middleware.UID(c)

	comment, err := h.svc.CreateComment(c.Request.Context(), uid, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPostMissing):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文章不存在"})
		case errors.Is(err, service.ErrCommentNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "父评论不存在"})
		case errors.Is(err, service.ErrParentMismatch):
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "父评论不属于当前文章"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建评论失败",
				"detail":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建评论成功",
		"data": gin.H{
			"id":         comment.Id,
			"content":    comment.Content,
			"user_id":    comment.UserId,
			"post_id":    comment.PostId,
			"parent_id":  comment.ParentId,
			"created_at": comment.CreatedAt.Format(time.RFC3339),
		},
	})
}

// DeleteComment 删除评论：仅评论作者本人可删除。
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	uid := middleware.UID(c)
	if err := h.svc.DeleteComment(c.Request.Context(), uid, uint(id)); err != nil {
		switch {
		case errors.Is(err, service.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "评论不存在"})
		case errors.Is(err, service.ErrCommentForbidden):
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该评论"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除评论失败",
				"detail":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除评论成功",
	})
}

// ListCommentsByPost 列出文章下的评论树。
func (h *CommentHandler) ListCommentsByPost(c *gin.Context) {
	postIdStr := c.Param("id")
	if postIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  "post_id不能为空",
		})
		return
	}
	postId, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	list, total, err := h.svc.ListCommentsByPost(c.Request.Context(), uint(postId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论失败",
			"detail":  err.Error(),
		})
		return
	}
	page, pageSize := 1, int(total)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询评论成功",
		"data": util.PageResult{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     list,
		},
	})
}

// ReplyComment 回复评论
func (h *CommentHandler) ReplyComment(c *gin.Context) {
	var req dto.ReplyCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}
	parentIdStr := c.Param("id")
	parentIdUint, err := strconv.ParseUint(parentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	uid := middleware.UID(c)
	comment, err := h.svc.ReplyToComment(c.Request.Context(), uid, uint(parentIdUint), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "父评论不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建回复失败",
				"detail":  err.Error(),
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建回复成功",
		"data": gin.H{
			"id":         comment.Id,
			"content":    comment.Content,
			"user_id":    comment.UserId,
			"post_id":    comment.PostId,
			"parent_id":  comment.ParentId,
			"created_at": comment.CreatedAt.Format(time.RFC3339),
		},
	})
}
