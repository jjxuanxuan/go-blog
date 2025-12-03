package service

import (
	"context"
	"errors"
	"strconv"

	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
	"go-blog/internal/util"
	"gorm.io/gorm"
)

// 认证相关错误定义。
var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
)

// AuthService 处理注册、登录和令牌刷新逻辑。
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 构造认证服务。
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register 注册新用户，包含重名校验与密码哈希。
func (s *AuthService) Register(ctx context.Context, req dto.CreateUserReq) (*model.User, error) {
	count, err := s.userRepo.CountByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserAlreadyExists
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		Role:     "user",
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 校验用户名密码并签发访问令牌与刷新令牌。
func (s *AuthService) Login(ctx context.Context, req dto.LoginReq) (string, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	if !util.CheckPassword(user.Password, req.Password) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := util.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := util.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// RefreshAccessToken 根据刷新令牌颁发新的访问令牌。
func (s *AuthService) RefreshAccessToken(_ context.Context, refreshToken string) (string, error) {
	claims, err := util.ParseToken(refreshToken)
	if err != nil {
		return "", ErrInvalidRefresh
	}
	uid64, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil || uid64 == 0 {
		return "", ErrInvalidRefresh
	}
	accessToken, err := util.GenerateAccessToken(uint(uid64), claims.Role)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
