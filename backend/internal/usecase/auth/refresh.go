package auth

import (
	"blog-aws-backend/internal/domain/session"
	"blog-aws-backend/internal/domain/user"
	"blog-aws-backend/internal/infrastructure/token"
	"context"
	"errors"
	"time"
)

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshUsecase struct {
	userRepository    user.Repository
	sessionRepository session.Repository
	tokenGenerator    TokenGenerator
}

func NewRefreshUsecase(userRepository user.Repository, sessionRepository session.Repository, tokenGenerator TokenGenerator) *RefreshUsecase {
	return &RefreshUsecase{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		tokenGenerator:    tokenGenerator,
	}
}

func (u *RefreshUsecase) Execute(ctx context.Context, req RefreshRequest) (RefreshResponse, error) {
	refreshTokenHash := token.HashRefreshToken(req.RefreshToken)

	foundSession, err := u.sessionRepository.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if errors.Is(err, session.ErrorNotFound) {
		return RefreshResponse{}, ErrorInvalidCredentials
	}
	if err != nil {
		return RefreshResponse{}, err
	}

	if time.Now().After(foundSession.ExpiresAt) {
		if err := u.sessionRepository.Revoke(ctx, foundSession.ID); err != nil {
			// log error
		}
		return RefreshResponse{}, ErrorInvalidCredentials
	}

	foundUser, err := u.userRepository.GetUserByID(ctx, foundSession.UserID)
	if err != nil {
		return RefreshResponse{}, err
	}

	accessToken, err := u.tokenGenerator.GenerateToken(foundSession.UserID, foundUser.Email, string(foundUser.Role))
	if err != nil {
		return RefreshResponse{}, err
	}

	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return RefreshResponse{}, err
	}

	newRefreshTokenHash := token.HashRefreshToken(newRefreshToken)

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	err = u.sessionRepository.Rotate(
		ctx,
		foundSession.ID,
		newRefreshTokenHash,
		expiresAt,
	)
	if err != nil {
		return RefreshResponse{}, err
	}

	return RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
