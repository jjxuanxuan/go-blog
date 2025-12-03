package service

import (
	"context"
	"errors"

	"go-blog/internal/dto"
	"go-blog/internal/model"
	"go-blog/internal/repository"
	"gorm.io/gorm"
)

// 评论相关错误定义。
var (
	ErrCommentNotFound  = errors.New("comment not found")
	ErrCommentForbidden = errors.New("comment forbidden")
	ErrParentMismatch   = errors.New("parent comment mismatch")
	ErrPostMissing      = errors.New("post not found")
)

// CommentService 聚合评论相关的业务逻辑。
type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

// NewCommentService 构造评论服务，注入评论与文章仓库。
func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment 创建评论，支持父子关系校验。
func (s *CommentService) CreateComment(ctx context.Context, uid uint, req dto.CreateCommentReq) (*model.Comment, error) {
	if _, err := s.postRepo.FindByID(ctx, req.PostId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostMissing
		}
		return nil, err
	}

	if req.ParentId != nil {
		parent, err := s.commentRepo.FindByID(ctx, *req.ParentId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCommentNotFound
			}
			return nil, err
		}
		if parent.PostId != req.PostId {
			return nil, ErrParentMismatch
		}
	}

	comment := &model.Comment{
		PostId:   req.PostId,
		UserId:   uid,
		Content:  req.Content,
		ParentId: req.ParentId,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// ReplyToComment 针对父评论创建回复，自动继承文章ID。
func (s *CommentService) ReplyToComment(ctx context.Context, uid, parentID uint, req dto.ReplyCommentReq) (*model.Comment, error) {
	parent, err := s.commentRepo.FindByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}

	comment := &model.Comment{
		PostId:   parent.PostId,
		UserId:   uid,
		Content:  req.Content,
		ParentId: &parent.Id,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment 删除评论，仅作者本人可操作。
func (s *CommentService) DeleteComment(ctx context.Context, uid, id uint) error {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}
	if comment.UserId != uid {
		return ErrCommentForbidden
	}
	return s.commentRepo.Delete(ctx, comment)
}

// ListCommentsByPost 根据文章构建评论树，并返回总数。
func (s *CommentService) ListCommentsByPost(ctx context.Context, postID uint) ([]dto.CommentResp, int64, error) {
	comments, err := s.commentRepo.ListByPostID(ctx, postID)
	if err != nil {
		return nil, 0, err
	}
	resp := buildCommentTree(comments)
	return resp, int64(len(comments)), nil
}

func buildCommentTree(list []model.Comment) []dto.CommentResp {
	m := make(map[uint]*dto.CommentResp, len(list))
	for _, c := range list {
		m[c.Id] = &dto.CommentResp{
			Id:       c.Id,
			Content:  c.Content,
			User:     dto.UserBrief{Id: c.User.ID, Username: c.User.Username},
			ParentId: c.ParentId,
			PostId:   c.PostId,
			Replies:  []dto.CommentResp{},
		}
	}

	for i := len(list) - 1; i >= 0; i-- {
		c := list[i]
		if c.ParentId != nil && *c.ParentId != c.Id {
			if parent, ok := m[*c.ParentId]; ok {
				parent.Replies = append(parent.Replies, *m[c.Id])
			}
		}
	}

	var roots []dto.CommentResp
	for _, c := range list {
		if c.ParentId == nil || *c.ParentId == c.Id {
			roots = append(roots, *m[c.Id])
			continue
		}
		if _, ok := m[*c.ParentId]; !ok {
			roots = append(roots, *m[c.Id])
		}
	}
	return roots
}
