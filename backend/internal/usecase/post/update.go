package post

import (
  postEntity "blog-aws-backend/internal/domain/post"
	"context"
	"time"
)

type UpdatePostRequest struct {
  ID string
	Title    string
	AuthorID string
	Content  string
}

type PostUpdateUsecase struct {
	postRepository postEntity.Repository
}

func NewPostUpdateUsecase(postRepository postEntity.Repository) *PostUpdateUsecase {
	return &PostUpdateUsecase{
		postRepository: postRepository,
	}
}

func (u *PostUpdateUsecase) Execute(ctx context.Context, req UpdatePostRequest) (postEntity.Post, error) {
  post := postEntity.Post{
    ID:           req.ID,
    Title:        req.Title,
    AuthorID:     req.AuthorID,
    Content:      req.Content,
    UpdatedAt:    time.Now(),
  }

  if err := u.postRepository.Update(ctx, post); err != nil {
    return postEntity.Post{}, err
  }

  return post, nil
}
