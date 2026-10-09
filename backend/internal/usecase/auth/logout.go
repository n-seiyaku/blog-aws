package auth

import (
	"blog-aws-backend/internal/domain/session"
	"blog-aws-backend/internal/infrastructure/token"
	"context"
	"errors"
)

type LogoutRequest struct {
	RefreshToken string
}

type LogoutUsecase struct {
	sessionRepository session.Repository
}

func NewLogoutUsecase(sessionRepository session.Repository) *LogoutUsecase {
	return &LogoutUsecase{
		sessionRepository: sessionRepository,
	}
}

func (u *LogoutUsecase) Execute(ctx context.Context, req LogoutRequest) error {
	refreshTokenHash := token.HashRefreshToken(req.RefreshToken)

	foundSession, err := u.sessionRepository.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if errors.Is(err, session.ErrorNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	return u.sessionRepository.Revoke(ctx, foundSession.ID)
}
