package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-blog/internal/service"
)

const (
	maxRequestSizeSingle = 10 << 20
	maxRequestSizeMulti  = 50 << 20
)

// UploadHandler 处理上传相关 HTTP 请求。
type UploadHandler struct{ svc *service.UploadService }

func NewUploadHandler(svc *service.UploadService) *UploadHandler { return &UploadHandler{svc: svc} }

// UploadSingle 单文件上传接口：POST /api/upload
func (h *UploadHandler) UploadSingle(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSizeSingle)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择要上传的文件",
			"detail":  err.Error(),
		})
		return
	}

	url, err := h.svc.UploadSingle(c.Request.Context(), file)
	if err != nil {
		h.renderUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data":    gin.H{"url": url},
	})
}

// UploadMulti 多文件上传接口：POST /api/upload/multi
func (h *UploadHandler) UploadMulti(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSizeMulti)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "解析表单失败",
			"detail":  err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "至少上传一个文件",
		})
		return
	}

	urls, err := h.svc.UploadMulti(c.Request.Context(), files)
	if err != nil {
		h.renderUploadError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data":    gin.H{"urls": urls},
	})
}

func (h *UploadHandler) renderUploadError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrInvalidUpload) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    500,
		"message": "文件上传失败",
		"detail":  err.Error(),
	})
}
