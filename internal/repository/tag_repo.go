package repository

import (
	"context"

	"go-blog/internal/model"
	"gorm.io/gorm"
)

// TagRepository 负责标签的存取。
type TagRepository struct {
	DB *gorm.DB
}

// NewTagRepository 创建标签仓库。
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// ListOrdered 按权重倒序返回标签列表。
func (r *TagRepository) ListOrdered(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.DB.WithContext(ctx).
		Order("weight DESC, id ASC").
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// Create 新增标签。
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	return r.DB.WithContext(ctx).Create(tag).Error
}
