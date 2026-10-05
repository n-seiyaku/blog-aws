package auth

import (
	"blog-aws-backend/internal/domain/session"
	"blog-aws-backend/internal/domain/user"
	"blog-aws-backend/internal/infrastructure/token"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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
	AccessToken  string
	RefreshToken string
	User         user.User
}

type LoginUsecase struct {
	userRepository    user.Repository
	sessionRepository session.Repository
	tokenGenerator    TokenGenerator
}

func NewLoginUsecase(
	userRepository user.Repository,
	sessionRepository session.Repository,
	tokenGenerator TokenGenerator,
) *LoginUsecase {
	return &LoginUsecase{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		tokenGenerator:    tokenGenerator,
	}
}

// func (u *LoginUsecase) Execute(
// 	ctx context.Context,
// 	req LoginRequest,
// ) (LoginResponse, error) {
// 	mockUser := user.User{
// 		ID:           "1",
// 		Email:        "test@gmail.com",
// 		PasswordHash: "$2a$10$cpD/X895o9Kj4tI8GjR0Nu9SIYe19iHuc7Fk4mOfJEaj1vi5okRqy",
// 		Role:         user.RoleUser,
// 	}

// 	email := strings.ToLower(strings.TrimSpace(req.Email))

// 	if email != mockUser.Email {
// 		return LoginResponse{}, ErrorInvalidCredentials
// 	}

// 	err := bcrypt.CompareHashAndPassword(
// 		[]byte(mockUser.PasswordHash),
// 		[]byte(req.Password),
// 	)
// 	if err != nil {
// 		return LoginResponse{}, ErrorInvalidCredentials
// 	}

// 	accessToken, err := u.tokenGenerator.GenerateToken(
// 		mockUser.ID,
// 		mockUser.Email,
// 		string(mockUser.Role),
// 	)
// 	if err != nil {
// 		return LoginResponse{}, err
// 	}

// 	return LoginResponse{
// 		AccessToken: accessToken,
// 		User:        mockUser,
// 	}, nil
// }

// to use database

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

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return LoginResponse{}, err
	}

	refreshTokenHash := token.HashRefreshToken(refreshToken)

	now := time.Now()
	expireAt := now.Add(30 * 24 * time.Hour)

	userID, err := uuid.NewV7()
	if err != nil {
		return LoginResponse{}, err
	}

	newSession := session.Session{
		ID:               userID.String(),
		UserID:           foundUser.ID,
		RefreshTokenHash: refreshTokenHash,
		CreatedAt:        now,
		ExpiresAt:        expireAt,
	}

	if err := u.sessionRepository.Create(ctx, newSession); err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         foundUser,
	}, nil
}
