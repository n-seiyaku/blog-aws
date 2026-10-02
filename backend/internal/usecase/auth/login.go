package auth

import (
	"blog-aws-backend/internal/domain/user"
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrorInvalidCredentials = errors.New("invalid credentials")

type TokenGenerator interface {
	GenerateToken(userID string, email string, role string) (string, error)
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken string
	User        user.User
}

type LoginUsecase struct {
	userRepository user.Repository
	tokenGenerator TokenGenerator
}

func NewLoginUsecase(userRepository user.Repository, tokenGenerator TokenGenerator) *LoginUsecase {
	return &LoginUsecase{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (u *LoginUsecase) Execute(
	ctx context.Context,
	req LoginRequest,
) (LoginResponse, error) {
	mockUser := user.User{
		ID:           "1",
		Email:        "test@gmail.com",
		PasswordHash: "$2a$10$cpD/X895o9Kj4tI8GjR0Nu9SIYe19iHuc7Fk4mOfJEaj1vi5okRqy",
		Role:         user.RoleUser,
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email != mockUser.Email {
		return LoginResponse{}, ErrorInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(mockUser.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return LoginResponse{}, ErrorInvalidCredentials
	}

	accessToken, err := u.tokenGenerator.GenerateToken(
		mockUser.ID,
		mockUser.Email,
		string(mockUser.Role),
	)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		User:        mockUser,
	}, nil
}

/*
func (u *LoginUsecase) Execute(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	foundUser, err := u.userRepository.FindByEmail(ctx, email)

	if errors.Is(err, user.ErrorNotFound) {
		return LoginResponse{}, ErrorInvalidCredentials
	}

	if err != nil {
		return LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password))

	if err != nil {
		return LoginResponse{}, ErrorInvalidCredentials
	}

	accessToken, err := u.tokenGenerator.GenerateToken(foundUser.ID, foundUser.Email, string(foundUser.Role))

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken: accessToken,
		User:        foundUser,
	}, nil
}
*/
