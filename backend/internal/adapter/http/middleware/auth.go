package middleware

import (
	"blog-aws-backend/internal/infrastructure/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	MissingAccessToken         = "missing_access_token"
	InvalidAuthorizationHeader = "invalid_authorization_header"
	InvalidAccessToken         = "invalid_access_token"
	AccessTokenExpired         = "access_token_expired"
)

type AccessTokenVerifier interface {
	VerifyToken(tokenString string) (*token.Claims, error)
}

type AuthMiddleware struct {
	tokenVerifier AccessTokenVerifier
}

func NewAuthMiddleware(tokenVerifier AccessTokenVerifier) *AuthMiddleware {
	return &AuthMiddleware{tokenVerifier: tokenVerifier}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": MissingAccessToken,
			})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authorization, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": InvalidAuthorizationHeader,
			})
			return
		}

		accessToken := strings.TrimPrefix(authorization, prefix)

		claims, err := m.tokenVerifier.VerifyToken(accessToken)
		if err != nil {
			if err == token.ErrorAccessTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": AccessTokenExpired,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": InvalidAccessToken,
			})
			return
		}

		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Set("userID", claims.Subject)

		c.Next()
	}
}
