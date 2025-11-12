package handler

import (
	"errors"
	"go-blog/internal/dto"
	"go-blog/internal/middleware"
	"go-blog/internal/model"
	"go-blog/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentHandler struct {
	DB *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{DB: db}
}

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

	var post model.Post
	if err := h.DB.First(&post, req.PostId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询文章失败",
			"detail":  err.Error(),
		})
		return
	}

	comment := model.Comment{
		PostId:  req.PostId,
		UserId:  uid,
		Content: req.Content,
	}
	if err := h.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建评论失败",
			"detail":  err.Error(),
		})
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
			"created_at": comment.CreatedAt.Format(time.RFC3339),
		},
	})
}

// FindCommentById 获取评论详情（预加载用户），错误时已写入响应
func (h *CommentHandler) FindCommentById(c *gin.Context, id uint) (*model.Comment, bool) {
	var comment model.Comment
	if err := h.DB.Preload("User").First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "评论不存在",
			})
			return nil, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论失败",
			"detail":  err.Error(),
		})
		return nil, false
	}
	return &comment, true
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
	comment, ok := h.FindCommentById(c, uint(id))
	if !ok {
		return
	}

	if comment.UserId != uid {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权操作该评论",
		})
		return
	}

	if err := h.DB.Delete(comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除评论失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除评论成功",
	})
}

func (h *CommentHandler) ListCommentsByPost(c *gin.Context) {
	postIdStr := c.Param("post_id")
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
	page, pageSize := util.ParsePage(c)
	offset := (page - 1) * pageSize

	var total int64
	if err := h.DB.Model(&model.Comment{}).Where("post_id=?", uint(postId)).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论总数失败",
			"detail":  err.Error(),
		})
		return
	}
	var comments []model.Comment
	if err := h.DB.Where("post_id=?", uint(postId)).Preload("User").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询评论失败",
			"detail":  err.Error(),
		})
		return
	}
	respList := make([]dto.CommentResp, 0, len(comments))
	for _, comment := range comments {
		if comment.ParentId == nil {
			respList = append(respList, buildCommentResp(comment))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询评论成功",
		"data": util.PageResult{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     respList,
		},
	})
}

func buildCommentResp(cmt model.Comment) dto.CommentResp {
	// 先填自己这一层
	resp := dto.CommentResp{
		Id:      cmt.Id,
		Content: cmt.Content,
		User: dto.UserBrief{
			Id:       cmt.User.ID,
			Username: cmt.User.Username, // 假设 User 有 Username 字段
		},
		PostId: cmt.PostId,
	}

	// ParentId 是 *uint
	if cmt.ParentId != nil {
		resp.ParentId = cmt.ParentId
	}

	// 再递归构造子评论
	if len(cmt.Replies) > 0 {
		resp.Replies = make([]dto.CommentResp, 0, len(cmt.Replies))
		for _, r := range cmt.Replies {
			resp.Replies = append(resp.Replies, buildCommentResp(r))
		}
	}

	return resp
}
