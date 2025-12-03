package service

import (
	"context"

	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
)

// CategoryService 处理分类相关的业务逻辑。
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService 构造分类服务，注入仓库。
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// ListCategories 返回按排序字段排好的分类列表。
func (s *CategoryService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.ListOrdered(ctx)
}

// CreateCategory 创建新的分类。
func (s *CategoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryReq) (*model.Category, error) {
	category := &model.Category{
		Name:     req.Name,
		ParentId: req.ParentId,
		Slug:     req.Slug,
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}
