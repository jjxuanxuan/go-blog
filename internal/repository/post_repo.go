package repository

import (
	"context"
	"time"

	"go-blog/internal/model"
	"gorm.io/gorm"
)

// PostFilter 文章列表筛选条件。
type PostFilter struct {
	CategoryID *uint   //分类 id
	TagIDs     []uint  //标签 id 列表
	Status     *string //状态
	Keyword    string  //关键词：title/content 模糊查询
	Order      string  //latest / hot
	Page       int
	PageSize   int
}

// PostRepository 提供文章的存取与查询。
type PostRepository struct {
	DB *gorm.DB
}

// NewPostRepository 创建文章仓库。
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// WithDB 用于在事务中替换为 tx
func (r *PostRepository) WithDB(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// ListPosts 根据过滤条件分页查询文章。
func (r *PostRepository) ListPosts(ctx context.Context, f PostFilter) (posts []model.Post, total int64, err error) {
	db := r.DB.WithContext(ctx).Model(&model.Post{}).Preload("Category").Preload("Tags").Preload("User")
	if f.CategoryID != nil && *f.CategoryID > 0 {
		db = db.Where("category_id = ?", *f.CategoryID)
	}

	if f.Status != nil && *f.Status != "" {
		db = db.Where("status = ?", *f.Status)
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
// ListByUserID 查询某用户的文章列表。
func (r *PostRepository) ListByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	var posts []model.Post
	if err := r.DB.Where("user_id = ?", userID).
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// Create 创建文章
func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	return r.DB.WithContext(ctx).Create(post).Error
}

// FindByID 根据 ID 查询文章（不预加载）
func (r *PostRepository) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	if err := r.DB.WithContext(ctx).First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// Save 保存文章（更新）
func (r *PostRepository) Save(ctx context.Context, post *model.Post) error {
	return r.DB.WithContext(ctx).Save(post).Error
}

// Delete 删除文章
func (r *PostRepository) Delete(ctx context.Context, post *model.Post) error {
	return r.DB.WithContext(ctx).Delete(post).Error
}

// ReplaceTags 替换文章标签
func (r *PostRepository) ReplaceTags(ctx context.Context, post *model.Post, tagIDs []uint) error {
	var tags []model.Tag
	if len(tagIDs) > 0 {
		if err := r.DB.WithContext(ctx).
			Where("id IN ?", tagIDs).
			Find(&tags).Error; err != nil {
			return err
		}
	}
	// Association 暂时没 ctx 概念，这里直接用 Model(post)
	return r.DB.Model(post).Association("Tags").Replace(&tags)
}

// FindAllWithUser 查询所有文章并预加载作者
func (r *PostRepository) FindAllWithUser(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	if err := r.DB.WithContext(ctx).
		Preload("User").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// CountAll 统计文章总数。
func (r *PostRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.Post{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountSince 统计指定时间后的文章数。
func (r *PostRepository) CountSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&model.Post{}).
		Where("created_at >= ?", since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// TopPostStat 表示评论最多的文章统计。
type TopPostStat struct {
	PostID       uint
	Title        string
	CommentCount int64
}

// TopPostsByComments 按评论数降序获取热门文章。
func (r *PostRepository) TopPostsByComments(ctx context.Context, limit int) ([]TopPostStat, error) {
	if limit <= 0 {
		limit = 5
	}
	var stats []TopPostStat
	err := r.DB.WithContext(ctx).
		Table("posts").
		Select("posts.id as post_id, posts.title as title, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Where("posts.status = ?", "published").
		Group("posts.id").
		Order("comment_count DESC, posts.id DESC").
		Limit(limit).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}
