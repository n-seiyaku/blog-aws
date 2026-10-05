package session

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, session Session) error
	FindByRefreshTokenHash(ctx context.Context, refreshToken string) (Session, error)
	Rotate(ctx context.Context, sessionID string, refreshToken string, expiresAt time.Time) error
	Revoke(ctx context.Context, sessionID string) error
}
