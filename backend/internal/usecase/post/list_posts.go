package post

import (
	"blog-aws-backend/internal/domain/post"
	"context"
)

const (
	MinLimit = 10
	MaxLimit = 50
)

type ListPostsRequest struct {
	ID       string
	Title    string
	AuthorID string
	Content  string
}

type ListPostsUsecase struct {
	postRepository post.Repository
}

func NewListPostsUsecase(postRepository post.Repository) *ListPostsUsecase {
	return &ListPostsUsecase{
		postRepository: postRepository,
	}
}

func (u *ListPostsUsecase) Execute(ctx context.Context, limit int32, cursor string) ([]post.Post, string, error) {
	if limit <= 0 {
		limit = MinLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return u.postRepository.ListPosts(ctx, limit, cursor)
}
