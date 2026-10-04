package post

import (
	"blog-aws-backend/internal/domain/post"
	"context"
	"time"
)

type UpdatePostRequest struct {
	ID       string
	Title    string
	AuthorID string
	Content  string
}

type UpdatePostUsecase struct {
	postRepository post.Repository
}

func NewUpdatePostUsecase(postRepository post.Repository) *UpdatePostUsecase {
	return &UpdatePostUsecase{
		postRepository: postRepository,
	}
}

func (u *UpdatePostUsecase) Execute(ctx context.Context, req UpdatePostRequest) (post.Post, error) {
	updatedPost := post.Post{
		ID:        req.ID,
		Title:     req.Title,
		AuthorID:  req.AuthorID,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	if err := u.postRepository.Update(ctx, updatedPost); err != nil {
		return post.Post{}, err
	}

	return updatedPost, nil
}
