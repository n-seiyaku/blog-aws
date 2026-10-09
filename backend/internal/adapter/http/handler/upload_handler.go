package handler

import (
	"blog-aws-backend/internal/usecase/upload"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type presignUploadRequest struct {
	ContentType string `json:"contentType" binding:"required"`
	Size        int64  `json:"size" binding:"required"`
}

type UploadResponse struct {
	UploadURL string `json:"uploadUrl"`
	ImageKey  string `json:"imageKey"`
}

type UploadHandler struct {
	presign *upload.PresignUploadUseCase
}

func NewUploadHandler(presign *upload.PresignUploadUseCase) *UploadHandler {
	return &UploadHandler{
		presign: presign,
	}
}
func (h *UploadHandler) Presign(c *gin.Context) {
	var req presignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	fmt.Println("userId =======================================", c.GetString("userID"))
	res, err := h.presign.Execute(c.Request.Context(), upload.PresignUploadRequest{
		UserID:      c.GetString("userID"),
		ContentType: req.ContentType,
		Size:        req.Size,
	})
	if err != nil {
		if errors.Is(err, upload.ErrInvalidImage) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid image",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate upload url",
		})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		UploadURL: res.UploadURL,
		ImageKey:  res.ImageKey,
	})
}
