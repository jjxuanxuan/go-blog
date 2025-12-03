package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-blog/internal/dto"
	"go-blog/internal/service"
)

// CategoryHandler 处理分类相关 HTTP 请求。
type CategoryHandler struct{ svc *service.CategoryService }

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// ListCategories GET /api/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	cats, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询分类失败",
			"detail":  err.Error(),
		})
		return
	}
	resp := make([]dto.CategoryResp, 0, len(cats))
	for _, cat := range cats {
		resp = append(resp, dto.CategoryResp{
			Id:       cat.Id,
			Name:     cat.Name,
			Slug:     cat.Slug,
			ParentId: cat.ParentId,
		})
	}

	c.JSON(http.StatusOK, gin.H{
	"code": 0,
	"data": resp,
})
}

// CreateCategory 创建分类。
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"detail":  err.Error(),
		})
		return
	}

	if _, err := h.svc.CreateCategory(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建分类失败",
			"detail":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}
