package main

import (
	"blog-aws-backend/internal/adapter/http/handler"
	"blog-aws-backend/internal/adapter/http/middleware"
	"blog-aws-backend/internal/infrastructure/database"
	"blog-aws-backend/internal/infrastructure/storage"
	"blog-aws-backend/internal/infrastructure/token"
	"blog-aws-backend/internal/usecase/auth"
	"blog-aws-backend/internal/usecase/post"
	"blog-aws-backend/internal/usecase/upload"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	postTableName := os.Getenv("DYNAMODB_POSTS_TABLE")
	if postTableName == "" {
		log.Fatal("DYNAMODB_POSTS_TABLE is required")
	}

	s3Bucket := os.Getenv("S3_POST_IMAGES_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_POST_IMAGES_BUCKET is required")
	}

	// AWS config
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	s3Presigner := s3.NewPresignClient(s3Client)
	s3Storage := storage.NewS3Storage(s3Client, s3Presigner, s3Bucket)

	// Infrastructure
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	userRepository := database.NewUserRepository(
		dynamoClient,
		userTableName,
	)

	sessionRepository := database.NewSessionRepository(
		dynamoClient,
		sessionTableName,
	)

	postRepository := database.NewPostRepository(
		dynamoClient,
		postTableName,
	)

	// JWT
	jwtGenerator := token.NewJWTGenerator(jwtSecret, 15*time.Minute)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtGenerator)

	// Usecases
	// Auth usecases
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

	// Post usecases
	createPostUsecase := post.NewCreatePostUsecase(postRepository, s3Storage)
	listPostsUsecase := post.NewListPostsUsecase(postRepository)
	deletePostUsecase := post.NewDeletePostUsecase(postRepository)
	updatePostUsecase := post.NewUpdatePostUsecase(postRepository, s3Storage)

	// Upload usecases
	presignUploadUsecase := upload.NewPresignUploadUseCase(s3Storage)

	// HTTP handler
	authHandler := handler.NewAuthHandler(
		loginUsecase,
		logoutUsecase,
		refreshUsecase,
		registerUsecase,
	)

	postHandler := handler.NewPostHandler(
		createPostUsecase,
		updatePostUsecase,
		deletePostUsecase,
		listPostsUsecase,
	)

	uploadHandler := handler.NewUploadHandler(presignUploadUsecase)

	// Gin
	router := gin.Default()

	router.Use(middleware.CORS())

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

	protected := router.Group("/")
	protected.Use(authMiddleware.Handle())
	{
		protected.GET("/posts", postHandler.ListPosts)
		protected.POST("/posts", postHandler.Create)
		protected.PUT("/posts/:id", postHandler.Update)
		protected.DELETE("/posts/:id", postHandler.Delete)
		protected.POST("/uploads/presign", uploadHandler.Presign)
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
