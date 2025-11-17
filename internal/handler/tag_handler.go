package handler

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/dto"
	"go-blog/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type TagHandler struct {
	DB *gorm.DB
}

func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{DB: db}
}

func (h *TagHandler) ListTags(c *gin.Context) {
	var tags []model.Tag
	if err := h.DB.Order("weight DESC, id ASC").Find(&tags).Error; err != nil {
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

	tag := model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := h.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    400,
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
