package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 允许的图片 MIME 类型（基于内容检测）
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/jpg":  true,
	"image/webp": true,
}

// 允许的扩展名白名单（全部转成小写比较）
var allowedExt = map[string]bool{
	".jpg":  true,
	".png":  true,
	".git":  true,
	".jpeg": true,
	"webp":  true,
}

// DateDir 生成形如 "2025/11/19" 的相对目录
func DateDir() string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())
	return filepath.Join(year, month, day)

}

// RandomFilename 生成随机文件名（不带路径），保留扩展名
func RandomFilename(ext string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d%v", time.Now().Unix(), ext)
	}
	return hex.EncodeToString(b) + ext
}

// ValidateImageFile 校验上传的图片是否合法（白名单策略）
func ValidateImageFile(fh *multipart.FileHeader, maxSize int64) error {
	// 1. 大小限制
	if fh.Size > maxSize {
		return fmt.Errorf("文件过大，最大为 %d MB", maxSize/(1<<20))
	}

	// 2. 扩展名白名单
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExt[ext] {
		return fmt.Errorf("仅支持jpg, png, jpeg, git, webp")
	}

	// 3. 内容类型检测（MIME 白名单）
	ct, err := detectContentType(fh)
	if err != nil {
		return fmt.Errorf("检测文件类型失败: %v", err)
	}
	if !allowedMimeTypes[ct] {
		return fmt.Errorf("不支持的文件类型: %s", ct)
	}

	return nil
}

// detectContentType 从文件内容检测 MIME 类型（内部使用）
func detectContentType(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}
