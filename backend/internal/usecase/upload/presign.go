package upload

import (
	"blog-aws-backend/internal/shared/image"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

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
	ext, ok := image.Extension(req.ContentType)
	if !ok || req.Size <= 0 || req.Size > image.MaxSize {
		return PresignUploadResponse{}, ErrInvalidImage
	}

	if req.UserID == "" || strings.ContainsAny(req.UserID, `/\`) {
		return PresignUploadResponse{}, ErrInvalidImage
	}

	imageUUID, err := uuid.NewV7()
	if err != nil {
		return PresignUploadResponse{}, err
	}

	key := fmt.Sprintf(
		"posts/%s/%s%s",
		req.UserID,
		imageUUID,
		ext,
	)

	url, err := u.storage.PresignUpload(ctx, key, req.ContentType)
	if err != nil {
		return PresignUploadResponse{}, err
	}

	return PresignUploadResponse{
		UploadURL: url,
		ImageKey:  key,
	}, nil
}
