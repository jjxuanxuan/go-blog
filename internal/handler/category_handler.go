package handler

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/dto"
	"go-blog/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

// ListCategories GET /api/categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var cats []model.Category
	if err := h.DB.Order("sort ASC,id ASC").Find(&cats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询分类失败",
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

	cat := model.Category{
		Name:     req.Name,
		ParentId: req.ParentId,
		Slug:     req.Slug,
	}

	if err := h.DB.Create(&cat).Error; err != nil {
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
