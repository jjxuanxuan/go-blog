package service

import (
	"context"
	"errors"
	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
	"gorm.io/gorm"
)

// 文章业务相关错误定义。
var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
)

// PostService 负责文章相关的业务逻辑
type PostService struct {
	DB   *gorm.DB
	Repo *repository.PostRepository
}

// NewPostService 构造文章服务，注入数据库和仓库。
func NewPostService(db *gorm.DB, repo *repository.PostRepository) *PostService {
	return &PostService{
		DB:   db,
		Repo: repo,
	}
}

// CreatePost 创建文章（带标签，使用事务保证文章和标签绑定一致）
func (s *PostService) CreatePost(ctx context.Context, uid uint, req dto.CreatePostReq) (*model.Post, error) {
	post := &model.Post{
		Title:      req.Title,
		Content:    req.Content,
		CategoryId: req.CategoryId,
		Status:     req.Status,
		UserID:     uid,
	}

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.Repo.WithDB(tx)

		// 1. 写文章
		if err := repoTx.Create(ctx, post); err != nil {
			return err
		}

		// 2. 处理标签
		if len(req.TagIds) > 0 {
			if err := repoTx.ReplaceTags(ctx, post, req.TagIds); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

// GetAllPosts 获取所有文章（预加载作者）
func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	return s.Repo.FindAllWithUser(ctx)
}

// GetPostByID 根据 id 查询文章详情（预加载作者）
func (s *PostService) GetPostByID(ctx context.Context, id uint) (*model.Post, error) {
	// 这里直接用 repo 的基础查询
	post, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	// 如果你希望这里也预加载 User，可以：
	// 1) 把 FindByID 改成 Preload("User")
	// 2) 或者另写一个 FindByIDWithUser
	return post, nil
}

// UpdatePost 更新文章：仅作者本人可更新，空字段不覆盖，标签一起维护
func (s *PostService) UpdatePost(ctx context.Context, uid, id uint, req dto.UpdatePostReq) (*model.Post, error) {
	var post *model.Post

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.Repo.WithDB(tx)

		// 1. 查询文章
		p, err := repoTx.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPostNotFound
			}
			return err
		}
		post = p

		// 2. 鉴权：只能作者自己改
		if post.UserID != uid {
			return ErrForbidden
		}

		// 3. 按需更新字段
		if req.Title != nil {
			post.Title = *req.Title
		}
		if req.Content != nil {
			post.Content = *req.Content
		}
		if req.Status != nil {
			post.Status = *req.Status
		}
		if req.CategoryID != nil {
			post.CategoryId = *req.CategoryID
		}

		// 4. 保存文章
		if err := repoTx.Save(ctx, post); err != nil {
			return err
		}

		// 5. 标签
		if err := repoTx.ReplaceTags(ctx, post, req.TagIDs); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePost 删除文章：仅作者本人可删
func (s *PostService) DeletePost(ctx context.Context, uid, id uint) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.Repo.WithDB(tx)

		post, err := repoTx.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPostNotFound
			}
			return err
		}

		if post.UserID != uid {
			return ErrForbidden
		}

		if err := repoTx.Delete(ctx, post); err != nil {
			return err
		}

		return nil
	})
}

// ListPosts 列表查询：直接复用 Repo 的过滤逻辑
func (s *PostService) ListPosts(ctx context.Context, f repository.PostFilter) ([]model.Post, int64, error) {
	return s.Repo.ListPosts(ctx, f)
}
