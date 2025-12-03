package service

import (
	"context"
	"errors"
	"go-blog/internal/model"
	"go-blog/internal/repository"
)

// 用户业务错误定义。
var (
	ErrorForbidden = errors.New("forbidden")
)

// UserService 处理用户个人信息与文章列表业务。
type UserService struct {
	UserRepo *repository.UserRepository
	PostRepo *repository.PostRepository
}

// NewUserService 构造用户服务。
func NewUserService(userRepo *repository.UserRepository, postRepo *repository.PostRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		PostRepo: postRepo,
	}
}

// GetMe 查询当前用户信息。
func (s *UserService) GetMe(cxt context.Context, uid uint) (*model.User, error) {
	return s.UserRepo.FindByID(cxt, uid)
}

// ListUserPosts 返回指定用户的文章，需本人访问。
func (s *UserService) ListUserPosts(cxt context.Context, requesterID, targetUserID uint) ([]model.Post, error) {
	if requesterID != targetUserID {
		return nil, ErrorForbidden
	}
	return s.PostRepo.ListByUserID(cxt, targetUserID)
}
