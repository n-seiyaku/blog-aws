package user

import (
	"context"
	"errors"
)

var ErrorNotFound = errors.New("user not found")
var ErrorEmailAlreadyExists = errors.New("email already exists")

type Repository interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}
