package session

import (
	"errors"
	"time"
)

var ErrorNotFound = errors.New("session not found")

type Session struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	CreatedAt        time.Time
	ExpiresAt        time.Time
}
