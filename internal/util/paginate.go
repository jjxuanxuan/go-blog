package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageResult 通用分页返回结构体，包含页码、总数和列表数据。
type PageResult struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	List     any   `json:"list"`
}

// ParsePage 解析分页参数，提供默认值并做合法性兜底。
func ParsePage(c *gin.Context) (page, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("page_size", "10")
	page, _ = strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	pageSize, _ = strconv.Atoi(sizeStr)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return
}
