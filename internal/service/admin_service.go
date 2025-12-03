package service

import (
	"context"
	"time"

	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
)

// AdminService 负责后台仪表盘与各类列表查询。
type AdminService struct {
	userRepo    *repository.UserRepository
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
}

// NewAdminService 构造 AdminService，并注入所需仓库。
func NewAdminService(userRepo *repository.UserRepository, postRepo *repository.PostRepository, commentRepo *repository.CommentRepository) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// Dashboard 汇总近一周增长与Top文章，供后台仪表盘展示。
func (s *AdminService) Dashboard(ctx context.Context, topN int) (*dto.AdminDashboardResp, error) {
	since := time.Now().AddDate(0, 0, -7)

	userCount, err := s.userRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	postCount, err := s.postRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	commentCount, err := s.commentRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	newUsers, err := s.userRepo.CountSince(ctx, since)
	if err != nil {
		return nil, err
	}

	newPosts, err := s.postRepo.CountSince(ctx, since)
	if err != nil {
		return nil, err
	}

	newComments, err := s.commentRepo.CountSince(ctx, since)
	if err != nil {
		return nil, err
	}

	topPostStats, err := s.postRepo.TopPostsByComments(ctx, topN)
	if err != nil {
		return nil, err
	}

	topPosts := make([]dto.AdminTopPost, 0, len(topPostStats))
	for _, stat := range topPostStats {
		topPosts = append(topPosts, dto.AdminTopPost{
			PostID:       stat.PostID,
			Title:        stat.Title,
			CommentCount: stat.CommentCount,
		})
	}

	return &dto.AdminDashboardResp{
		Metrics: dto.AdminDashboardMetrics{
			Users:    userCount,
			Posts:    postCount,
			Comments: commentCount,
		},
		Recent: dto.AdminRecentStats{
			NewUsers:    newUsers,
			NewPosts:    newPosts,
			NewComments: newComments,
		},
		TopPosts: topPosts,
	}, nil
}

// ListUsers 按管理员查询条件分页返回用户列表。
func (s *AdminService) ListUsers(ctx context.Context, q dto.AdminUserQuery) ([]model.User, int64, error) {
	filter := repository.UserFilter{
		Keyword:  q.Keyword,
		Role:     q.Role,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	return s.userRepo.List(ctx, filter)
}

// ListPosts 按状态、关键词等条件分页返回文章列表。
func (s *AdminService) ListPosts(ctx context.Context, q dto.AdminPostQuery) ([]model.Post, int64, error) {
	filter := repository.PostFilter{
		Keyword:  q.Keyword,
		Status:   q.Status,
		Page:     q.Page,
		PageSize: q.PageSize,
		Order:    "latest",
	}
	return s.postRepo.ListPosts(ctx, filter)
}

// ListComments 按用户/文章和关键词过滤评论并分页返回。
func (s *AdminService) ListComments(ctx context.Context, q dto.AdminCommentQuery) ([]model.Comment, int64, error) {
	filter := repository.CommentFilter{
		UserID:   q.UserID,
		PostID:   q.PostID,
		Keyword:  q.Keyword,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	return s.commentRepo.List(ctx, filter)
}
