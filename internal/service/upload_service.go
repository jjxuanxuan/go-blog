package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"go-blog/internal/repository"
	"go-blog/internal/util"
)

const (
	staticUploadPath = "/static/uploads"
	maxUploadSize    = 5 << 20
)

// UploadService 处理文件上传存储逻辑。
type UploadService struct {
	repo *repository.UploadRepository
}

// ErrInvalidUpload 表示文件类型或大小不合法。
var ErrInvalidUpload = errors.New("invalid upload")

// NewUploadService 构造上传服务。
func NewUploadService(repo *repository.UploadRepository) *UploadService {
	return &UploadService{repo: repo}
}

// UploadSingle 处理单文件上传并返回可访问URL。
func (s *UploadService) UploadSingle(_ context.Context, file *multipart.FileHeader) (string, error) {
	return s.persistFile(file)
}

// UploadMulti 处理多文件上传并返回URL列表。
func (s *UploadService) UploadMulti(_ context.Context, files []*multipart.FileHeader) ([]string, error) {
	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, err := s.persistFile(file)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// persistFile 校验文件并保存到磁盘，返回相对路径。
func (s *UploadService) persistFile(file *multipart.FileHeader) (string, error) {
	if err := util.ValidateImageFile(file, maxUploadSize); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidUpload, err)
	}
	subDir := util.DateDir()
	ext := filepath.Ext(file.Filename)
	filename := util.RandomFilename(ext)
	relative, err := s.repo.SaveFile(subDir, filename, file)
	if err != nil {
		return "", err
	}
	return staticUploadPath + "/" + relative, nil
}
