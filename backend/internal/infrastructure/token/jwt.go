package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTGenerator(secret string, ttl time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

type claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (g *JWTGenerator) GenerateToken(userID string, email string, role string) (string, error) {
	now := time.Now()

	tokenClaims := claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token.SignedString(g.secret)
}
