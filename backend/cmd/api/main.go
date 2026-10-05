package main

import (
	"blog-aws-backend/internal/adapter/http/handler"
	"blog-aws-backend/internal/adapter/http/middleware"
	"blog-aws-backend/internal/infrastructure/database"
	"blog-aws-backend/internal/infrastructure/token"
	"blog-aws-backend/internal/usecase/auth"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	userTableName := os.Getenv("DYNAMODB_USERS_TABLE")
	if userTableName == "" {
		log.Fatal("DYNAMODB_USERS_TABLE is required")
	}

	sessionTableName := os.Getenv("DYNAMODB_SESSIONS_TABLE")
	if sessionTableName == "" {
		log.Fatal("DYNAMODB_SESSIONS_TABLE is required")
	}

	// AWS config
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// Infrastructure
	dynamoClient := dynamodb.NewFromConfig(cfg)

	userRepository := database.NewUserRepository(
		dynamoClient,
		userTableName,
	)

	sessionRepository := database.NewSessionRepository(
		dynamoClient,
		sessionTableName,
	)

	// JWT
	jwtGenerator := token.NewJWTGenerator(jwtSecret, 15*time.Minute)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtGenerator)

	// Usecases
	loginUsecase := auth.NewLoginUsecase(
		userRepository,
		sessionRepository,
		jwtGenerator,
	)

	registerUsecase := auth.NewRegisterUsecase(
		userRepository,
	)

	logoutUsecase := auth.NewLogoutUsecase(
		sessionRepository,
	)

	refreshUsecase := auth.NewRefreshUsecase(
		userRepository,
		sessionRepository,
		jwtGenerator,
	)

	// HTTP handler
	authHandler := handler.NewAuthHandler(
		loginUsecase,
		logoutUsecase,
		refreshUsecase,
		registerUsecase,
	)

	// Gin
	router := gin.Default()

	api := router.Group("/api")
	{
		auth := api.Group("auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/register", authHandler.Register)
		}
	}

	protect := router.Group("/")
	protect.Use(authMiddleware.Handle())
	{
		protect.GET("/posts", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "this is protect route",
				"role":    ctx.GetString("role"),
				"email":   ctx.GetString("email"),
				"userID":  ctx.GetString("userID"),
			})
		})
	}

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
