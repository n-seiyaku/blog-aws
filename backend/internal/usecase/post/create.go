package post

import (
	postEntity "blog-aws-backend/internal/domain/post"
	"context"
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title    string
	AuthorID string
	Content  string
}

type PostCreateUseCase struct {
	postRepository postEntity.Repository
}

func NewPostCreateUseCase(postRepository postEntity.Repository) *PostCreateUseCase {
	return &PostCreateUseCase{
		postRepository: postRepository,
	}
}

func (u *PostCreateUseCase) Execute(ctx context.Context, req CreatePostRequest) (postEntity.Post, error) {
	postID, err := uuid.NewV7()
	if err != nil {
		return postEntity.Post{}, err
	}

	post := postEntity.Post{
		ID:        postID.String(),
		Title:     req.Title,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.postRepository.Create(ctx, post); err != nil {
		return postEntity.Post{}, err
	}

	return post, nil
}
