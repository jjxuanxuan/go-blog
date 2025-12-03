package repository

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadRepository 负责文件写入磁盘。
type UploadRepository struct {
	Root string
}

// NewUploadRepository 创建上传仓库，root 为保存根目录。
func NewUploadRepository(root string) *UploadRepository {
	return &UploadRepository{Root: root}
}

// SaveFile 将上传的文件保存到指定子目录并返回相对路径。
func (r *UploadRepository) SaveFile(subDir, filename string, file *multipart.FileHeader) (string, error) {
	dstDir := filepath.Join(r.Root, subDir)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(dstDir, filename)
	if err := writeMultipartFile(file, dstPath); err != nil {
		return "", err
	}

	return filepath.ToSlash(filepath.Join(subDir, filename)), nil
}

// writeMultipartFile 将 multipart 文件流写入目标路径。
func writeMultipartFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return err
	}
	return nil
}
