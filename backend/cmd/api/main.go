package main

import (
	"blog-aws-backend/internal/adapter/http/handler"
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

	tableName := os.Getenv("DYNAMODB_USERS_TABLE")

	if tableName == "" {
		log.Fatal("DYNAMODB_USERS_TABLE is required")
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
		tableName,
	)

	jwtGenerator := token.NewJWTGenerator(jwtSecret, 24*time.Hour)

	// Usecases
	loginUsecase := auth.NewLoginUsecase(
		userRepository,
		jwtGenerator,
	)

	registerUsecase := auth.NewRegisterUsecase(
		userRepository,
	)

	// HTTP handler
	authHandler := handler.NewAuthHandler(
		loginUsecase,
		registerUsecase,
	)

	// Gin
	router := gin.Default()

	api := router.Group("/api")
	{
		auth := api.Group("auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
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
