package handler

import (
	"blog-aws-backend/internal/usecase/post"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageKey string `json:"imageKey"`
}

type UpdatePostRequest struct {
	Title    string  `json:"title" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	ImageKey *string `json:"imageKey"`
}

type PostResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageKey  string `json:"imageKey"`
	AuthorID  string `json:"authorId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PostHandler struct {
	createUsecase    *post.CreatePostUsecase
	updateUsecase    *post.UpdatePostUsecase
	deleteUsecase    *post.DeletePostUsecase
	listPostsUsecase *post.ListPostsUsecase
}

func NewPostHandler(
	createUsecase *post.CreatePostUsecase,
	updateUsecase *post.UpdatePostUsecase,
	deleteUsecase *post.DeletePostUsecase,
	listPostsUsecase *post.ListPostsUsecase,
) *PostHandler {
	return &PostHandler{
		createUsecase:    createUsecase,
		updateUsecase:    updateUsecase,
		deleteUsecase:    deleteUsecase,
		listPostsUsecase: listPostsUsecase,
	}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	userID := c.GetString("userID")
	newPost, err := h.createUsecase.Execute(c.Request.Context(), post.CreatePostRequest{
		Title:    req.Title,
		AuthorID: userID,
		Content:  req.Content,
		ImageKey: req.ImageKey,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create post",
		})
		return
	}

	c.JSON(http.StatusCreated, PostResponse{
		ID:        newPost.ID,
		Title:     newPost.Title,
		Content:   newPost.Content,
		ImageKey:  *newPost.ImageKey,
		AuthorID:  newPost.AuthorID,
		CreatedAt: newPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newPost.UpdatedAt.Format(time.RFC3339),
	})
}

func (h *PostHandler) Update(ctx *gin.Context) {
	var req UpdatePostRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedPost, err := h.updateUsecase.Execute(ctx.Request.Context(), post.UpdatePostRequest{
		ID:       ctx.Param("id"),
		Title:    req.Title,
		AuthorID: ctx.GetString("userID"),
		Content:  req.Content,
		ImageKey: req.ImageKey,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, PostResponse{
		ID:        updatedPost.ID,
		Title:     updatedPost.Title,
		Content:   updatedPost.Content,
		ImageKey:  *updatedPost.ImageKey,
		AuthorID:  updatedPost.AuthorID,
		CreatedAt: updatedPost.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updatedPost.UpdatedAt.Format(time.RFC3339),
	})
}

func (h *PostHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	err := h.deleteUsecase.Execute(c.Request.Context(), post.DeletePostRequest{
		ID: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete post",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	limitStr := c.Query("limit")
	cursor := c.Query("cursor")

	var limit int32 = 10
	if limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid limit",
			})
			return
		}
		limit = int32(parsedLimit)
	}

	posts, nextCursor, err := h.listPostsUsecase.Execute(c.Request.Context(), limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to list posts",
		})
		return
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			ImageKey:  *post.ImageKey,
			AuthorID:  post.AuthorID,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":      postResponses,
		"nextCursor": nextCursor,
	})
}
