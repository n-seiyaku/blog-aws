package post

import (
	"context"

	"blog-aws-backend/internal/domain/post"
)

type DeletePostRequest struct {
	ID string
}

type DeletePostUsecase struct {
	postRepository post.Repository
}

func NewDeletePostUsecase(postRepository post.Repository) *DeletePostUsecase {
	return &DeletePostUsecase{
		postRepository: postRepository,
	}
}

func (u *DeletePostUsecase) Execute(ctx context.Context, req DeletePostRequest) error {
	return u.postRepository.Delete(ctx, req.ID)
}
