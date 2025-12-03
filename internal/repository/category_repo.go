package repository

import (
	"context"

	"go-blog/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 提供分类的查询与写入。
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository 创建分类仓库。
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// ListOrdered 按排序字段获取分类列表。
func (r *CategoryRepository) ListOrdered(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	if err := r.DB.WithContext(ctx).
		Order("sort ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Create 新增分类。
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.DB.WithContext(ctx).Create(category).Error
}
