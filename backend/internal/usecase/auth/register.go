package auth

import (
	"blog-aws-backend/internal/domain/user"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string
	Password string
}

type RegisterUsecase struct {
	userRepository user.Repository
}

func NewRegisterUsecase(userRepository user.Repository) *RegisterUsecase {
	return &RegisterUsecase{
		userRepository: userRepository,
	}
}

func (u *RegisterUsecase) Execute(ctx context.Context, req RegisterRequest) (user.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// メールアドレスの重複チェック
	_, err := u.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return user.User{}, user.ErrorEmailAlreadyExists
	}

	if !errors.Is(err, user.ErrorNotFound) {
		return user.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return user.User{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return user.User{}, err
	}

	newUser := user.User{
		ID:           userID.String(),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         user.RoleUser,
		CreatedAt:    time.Now().UTC(),
	}

	if err := u.userRepository.Create(ctx, newUser); err != nil {
		return user.User{}, err
	}

	return newUser, nil
}
