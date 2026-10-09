package upload

import (
	"context"
	"errors"
	"strings"
)

const MaxImageSize int64 = 5 * 1024 * 1024

var ErrInvalidImage = errors.New("invalid image")

type Storage interface {
	PresignUpload(ctx context.Context,
		key string,
		contentType string,
	) (string, error)
}

type PresignUploadUseCase struct {
	storage Storage
}

type PresignUploadRequest struct {
	UserID      string
	ContentType string
	Size        int64
}

type PresignUploadResponse struct {
	UploadURL string
	ImageKey  string
}

func NewPresignUploadUseCase(storage Storage) *PresignUploadUseCase {
	return &PresignUploadUseCase{
		storage: storage,
	}
}

func (u *PresignUploadUseCase) Execute(
	ctx context.Context,
	req PresignUploadRequest,
) (PresignUploadResponse, error) {
	extension := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}

	ext, ok := extension[req.ContentType]
	if !ok || req.Size <= 0 || req.Size > MaxImageSize {
		return PresignUploadResponse{}, ErrInvalidImage
	}

	if req.UserID == "" || strings.ContainsAny(req.UserID, `/\`) {
		return PresignUploadResponse{}, ErrInvalidImage
	}

}
