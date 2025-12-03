package repository

import (
	"context"
	"time"

	"go-blog/internal/model"
	"gorm.io/gorm"
)

// UserRepository 封装用户表的查询与写入。
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository 创建用户仓库。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// FindByID 按ID查询用户。
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := r.DB.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUsername 按用户名查询用户。
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	if err := r.DB.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// CountByUsernameOrEmail 统计用户名或邮箱是否已存在。
func (r *UserRepository) CountByUsernameOrEmail(ctx context.Context, username, email string) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Create 创建用户。
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

// UserFilter 用户列表筛选条件。
type UserFilter struct {
	Keyword  string
	Role     string
	Page     int
	PageSize int
}

// CountAll 统计用户总数。
func (r *UserRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountSince 统计指定时间后的新增用户数。
func (r *UserRepository) CountSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("created_at >= ?", since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// List 按条件分页查询用户。
func (r *UserRepository) List(ctx context.Context, f UserFilter) ([]model.User, int64, error) {
	db := r.DB.WithContext(ctx).Model(&model.User{})

	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ?", like, like)
	}

	if f.Role != "" {
		db = db.Where("role = ?", f.Role)
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

	var users []model.User
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
