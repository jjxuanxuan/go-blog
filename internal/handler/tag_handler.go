package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-blog/internal/dto"
	"go-blog/internal/service"
)

// TagHandler 处理标签相关 HTTP 请求。
type TagHandler struct{ svc *service.TagService }

func NewTagHandler(svc *service.TagService) *TagHandler { return &TagHandler{svc: svc} }

// ListTags 获取标签列表。
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.svc.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询标签失败",
			"detail":  err.Error(),
		})
		return
	}
	resp := make([]dto.TagResp, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, dto.TagResp{
			Id:     t.Id,
			Name:   t.Name,
			Slug:   t.Slug,
			Weight: t.Weight,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    resp,
	})
}

// CreateTag 新增标签。
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req dto.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}
	if _, err := h.svc.CreateTag(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建标签失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
