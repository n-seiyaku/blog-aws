package post

import (
  postEntity "blog-aws-backend/internal/domain/post"
	"context"
	"time"
	"uuid"
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
  post := postEntity.Post{
    ID:           uuid.NewV7().String(),
    Title:        req.Title,
    AuthorID:     req.AuthorID,
    Content:      req.Content,
    CreatedAt:    time.Now(),
    UpdatedAt:    time.Now(),
  }

  if err := u.postRepository.Create(ctx, post); err != nil {
    return postEntity.Post{}, err
  }

  return post, nil
}
