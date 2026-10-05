package post

import (
	"context"

	"blog-aws-backend/internal/domain/post"
)

type PostDeleteRequest struct {
	ID string
}

type PostDeleteUsecase struct {
	postRepository post.Repository
}

func NewPostDeleteUsecase(postRepository post.Repository) *PostDeleteUsecase {
	return &PostDeleteUsecase{
		postRepository: postRepository,
	}
}

func (u *PostDeleteUsecase) Execute(ctx context.Context, req PostDeleteRequest) error {
	return u.postRepository.Delete(ctx, req.ID)
}
