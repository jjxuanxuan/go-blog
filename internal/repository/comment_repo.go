package repository

import (
	"context"
	"time"

	"go-blog/internal/model"
	"gorm.io/gorm"
)

// CommentRepository 提供评论表的CRUD与查询。
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository 创建评论仓库。
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// Create 新增评论。
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.DB.WithContext(ctx).Create(comment).Error
}

// FindByID 按ID查询评论。
func (r *CommentRepository) FindByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.DB.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete 删除评论。
func (r *CommentRepository) Delete(ctx context.Context, comment *model.Comment) error {
	return r.DB.WithContext(ctx).Delete(comment).Error
}

// ListByPostID 查询文章下的所有评论并预加载用户。
func (r *CommentRepository) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.DB.WithContext(ctx).
		Where("post_id = ?", postID).
		Preload("User").
		Order("id ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.Comment{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CommentRepository) CountSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&model.Comment{}).
		Where("created_at >= ?", since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CommentFilter 评论筛选条件。
type CommentFilter struct {
	UserID   *uint
	PostID   *uint
	Keyword  string
	Page     int
	PageSize int
}

// List 按条件分页查询评论列表。
func (r *CommentRepository) List(ctx context.Context, f CommentFilter) ([]model.Comment, int64, error) {
	db := r.DB.WithContext(ctx).Model(&model.Comment{}).Preload("User").Preload("Post")

	if f.UserID != nil && *f.UserID > 0 {
		db = db.Where("user_id = ?", *f.UserID)
	}

	if f.PostID != nil && *f.PostID > 0 {
		db = db.Where("post_id = ?", *f.PostID)
	}

	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		db = db.Where("content LIKE ?", like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := f.Page
	if page <= 0 {
		page = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var comments []model.Comment
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}
