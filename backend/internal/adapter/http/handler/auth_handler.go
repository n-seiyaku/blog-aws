package handler

import (
	"blog-aws-backend/internal/domain/user"
	"blog-aws-backend/internal/usecase/auth"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

type AuthHandler struct {
	loginUsecase    *auth.LoginUsecase
	logoutUsecase   *auth.LogoutUsecase
	refreshUsecase  *auth.RefreshUsecase
	registerUsecase *auth.RegisterUsecase
}

func NewAuthHandler(
	loginUsecase *auth.LoginUsecase,
	logoutUsecase *auth.LogoutUsecase,
	refreshUsecase *auth.RefreshUsecase,
	registerUsecase *auth.RegisterUsecase,
) *AuthHandler {
	return &AuthHandler{
		loginUsecase:    loginUsecase,
		logoutUsecase:   logoutUsecase,
		refreshUsecase:  refreshUsecase,
		registerUsecase: registerUsecase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	result, err := h.loginUsecase.Execute(c.Request.Context(), auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, auth.ErrorInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
		"user": gin.H{
			"id":    result.User.ID,
			"email": result.User.Email,
			"role":  result.User.Role,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	newUser, err := h.registerUsecase.Execute(c.Request.Context(), auth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, user.ErrorEmailAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email already exists",
		})
		return
	}

	if err != nil {
		// エラー詳細をログに出力
		log.Printf("Register failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to register",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    newUser.ID,
			"email": newUser.Email,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid request",
		})
		return
	}

	result, err := h.refreshUsecase.Execute(
		c.Request.Context(),
		auth.RefreshRequest{
			RefreshToken: req.RefreshToken,
		},
	)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	err := h.logoutUsecase.Execute(
		c.Request.Context(),
		auth.LogoutRequest{
			RefreshToken: req.RefreshToken,
		},
	)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
