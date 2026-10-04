package post

import (
	"context"
	"errors"
)

var ErrorNotFound = errors.New("post not found")

type Repository interface {
	Create(ctx context.Context, post Post) error
	Update(ctx context.Context, post Post) error
	Delete(ctx context.Context, id string) error
	ListPosts(ctx context.Context, limit int32, cursor string) ([]Post, string, error)
}
