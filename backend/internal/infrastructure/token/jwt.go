package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrorAccessTokenExpired = errors.New("access token expired")
	ErrorInvalidAccessToken = errors.New("invalid access token")
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

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (g *JWTGenerator) GenerateToken(userID string, email string, role string) (string, error) {
	now := time.Now()

	tokenClaims := Claims{
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

func (g *JWTGenerator) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return g.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrorAccessTokenExpired
		}

		return nil, ErrorInvalidAccessToken
	}

	if !token.Valid {
		return nil, ErrorInvalidAccessToken
	}

	return claims, nil
}
