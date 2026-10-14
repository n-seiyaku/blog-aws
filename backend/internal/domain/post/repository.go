package post

import (
	"context"
	"errors"
)

var (
	ErrNotFound     = errors.New("post not found")
	ErrForbidden    = errors.New("not allowed to update this post")
	ErrConflict     = errors.New("post was modified concurrently")
	ErrInvalidImage = errors.New("invalid post image")
)

type Repository interface {
	Create(ctx context.Context, post Post) error
	Update(ctx context.Context, p Post, expectedImageKey string) error
	Delete(ctx context.Context, id string) error
	ListPosts(ctx context.Context, limit int32, cursor string) ([]Post, string, error)
	FindById(ctx context.Context, id string) (Post, error)
}
