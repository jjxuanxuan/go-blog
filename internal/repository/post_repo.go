package repository

import (
	"context"
	"go-blog/internal/model"
	"gorm.io/gorm"
)

type PostFilter struct {
	CategoryID *uint   //分类 id
	TagIDs     []uint  //标签 id 列表
	Status     *string //状态
	Keyword    string  //关键词：title/content 模糊查询
	Order      string  //latest / hot
	Page       int
	PageSize   int
}

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) ListPosts(ctx context.Context, f PostFilter) (posts []model.Post, total int64, err error) {
	db := r.DB.WithContext(ctx).Model(&model.Post{}).Preload("Category").Preload("Tags")
	if f.CategoryID != nil && *f.CategoryID > 0 {
		db = db.Where("category_id = ?", *f.CategoryID)
	}

	if f.Status != nil && *f.Status != "" {
		db = db.Where("status = ?", f.Status)
	}

	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		db = db.Where("title like ? or content like ?", kw, kw)
	}

	if len(f.TagIDs) > 0 {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id IN ?", f.TagIDs).
			Group("posts.id")
	}

	switch f.Order {
	case "latest":
		db = db.Order("posts.created_at DESC")
	case "hot":
		db = db.Order("posts.id DESC")
	default:
		db = db.Order("posts.id DESC")
	}

	// 先统计总数
	if err = db.Count(&total).Error; err != nil {
		return
	}

	// 分页
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 || f.PageSize > 100 {
		f.PageSize = 10
	}
	offset := (f.Page - 1) * f.PageSize

	err = db.Offset(offset).Limit(f.PageSize).Find(&posts).Error
	return
}
