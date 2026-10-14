package post

import (
	"blog-aws-backend/internal/domain/post"
	"context"
	"log/slog"
	"time"
)

type UpdatePostRequest struct {
	ID       string
	Title    string
	AuthorID string
	Content  string
	ImageKey *string
}

type UpdatePostImageStorage interface {
	VerifyPostImage(ctx context.Context, userID string, key string) error
	Delete(ctx context.Context, key string) error
}

type UpdatePostUsecase struct {
	postRepository post.Repository
	imageStorage   UpdatePostImageStorage
}

func NewUpdatePostUsecase(
	postRepository post.Repository,
	imageStorage UpdatePostImageStorage,
) *UpdatePostUsecase {
	return &UpdatePostUsecase{
		postRepository: postRepository,
		imageStorage:   imageStorage,
	}
}

func (u *UpdatePostUsecase) Execute(ctx context.Context, req UpdatePostRequest) (post.Post, error) {
	currentPost, err := u.postRepository.FindById(ctx, req.ID)
	if err != nil {
		return post.Post{}, err
	}

	updatedPost := currentPost
	updatedPost.Title = req.Title
	updatedPost.Content = req.Content
	updatedPost.UpdatedAt = time.Now().UTC()

	imageChanged := false

	currentImageKey := ""
	if currentPost.ImageKey != nil {
		currentImageKey = *currentPost.ImageKey
	}

	if req.ImageKey != nil {
		newImageKey := *req.ImageKey

		if newImageKey != currentImageKey {
			if newImageKey != "" {
				if err := u.imageStorage.VerifyPostImage(
					ctx,
					req.AuthorID,
					newImageKey,
				); err != nil {
					return post.Post{}, post.ErrInvalidImage
				}
			}

			updatedPost.ImageKey = &newImageKey
			imageChanged = true
		}
	}

	if err := u.postRepository.Update(ctx, updatedPost, currentImageKey); err != nil {
		return post.Post{}, err
	}

	if imageChanged && currentImageKey != "" {
		cleanupCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			5*time.Second,
		)
		defer cancel()

		if err := u.imageStorage.Delete(
			cleanupCtx,
			currentImageKey,
		); err != nil {
			slog.Warn(
				"failed to delete old post image",
				"postID", currentPost.ID,
				"imageKey", currentImageKey,
				"error", err,
			)
		}
	}
	return updatedPost, nil
}
