package service

import (
	"context"

	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
)

// TagService 处理标签相关业务。
type TagService struct {
	repo *repository.TagRepository
}

// NewTagService 构造标签服务。
func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

// ListTags 返回按权重排序的标签列表。
func (s *TagService) ListTags(ctx context.Context) ([]model.Tag, error) {
	return s.repo.ListOrdered(ctx)
}

// CreateTag 创建新标签。
func (s *TagService) CreateTag(ctx context.Context, req dto.CreateTagReq) (*model.Tag, error) {
	tag := &model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}
