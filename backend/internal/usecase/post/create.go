package post

import (
	"blog-aws-backend/internal/domain/post"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title    string
	AuthorID string
	Content  string
	ImageKey string
}

type PostImageStorage interface {
	VerifyPostImage(
		ctx context.Context,
		userID string,
		key string,
	) error
}

type CreatePostUsecase struct {
	postRepository post.Repository
	imageStorage   PostImageStorage
}

func NewCreatePostUsecase(postRepository post.Repository, imageStorage PostImageStorage) *CreatePostUsecase {
	return &CreatePostUsecase{
		postRepository: postRepository,
		imageStorage:   imageStorage,
	}
}

func (u *CreatePostUsecase) Execute(ctx context.Context, req CreatePostRequest) (post.Post, error) {
	postID, err := uuid.NewV7()
	if err != nil {
		return post.Post{}, err
	}

	now := time.Now()

	newPost := post.Post{
		ID:        postID.String(),
		Title:     req.Title,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		ImageKey:  req.ImageKey,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.ImageKey != "" {
		if err := u.imageStorage.VerifyPostImage(
			ctx,
			req.AuthorID,
			req.ImageKey,
		); err != nil {
			return post.Post{}, err
		}
	}

	if err := u.postRepository.Create(ctx, newPost); err != nil {
		return post.Post{}, err
	}

	return newPost, nil
}
