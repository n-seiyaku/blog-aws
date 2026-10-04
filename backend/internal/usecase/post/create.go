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
}

type CreatePostUsecase struct {
	postRepository post.Repository
}

func NewCreatePostUsecase(postRepository post.Repository) *CreatePostUsecase {
	return &CreatePostUsecase{
		postRepository: postRepository,
	}
}

func (u *CreatePostUsecase) Execute(ctx context.Context, req CreatePostRequest) (post.Post, error) {
	postID, err := uuid.NewV7()
	if err != nil {
		return post.Post{}, err
	}

	newPost := post.Post{
		ID:        postID.String(),
		Title:     req.Title,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.postRepository.Create(ctx, newPost); err != nil {
		return post.Post{}, err
	}

	return newPost, nil
}
