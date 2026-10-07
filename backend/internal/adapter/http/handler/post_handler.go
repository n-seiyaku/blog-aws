package handler

import (
	"blog-aws-backend/internal/usecase/post"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	AuthorID string `json:"authorId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	ID       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	AuthorID string `json:"authorId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type DeletePostRequest struct {
	ID string `json:"id" binding:"required"`
}

type PostHandler struct {
	createUsecase    *post.CreatePostUsecase
	updateUsecase    *post.UpdatePostUsecase
	deleteUsecase    *post.DeletePostUsecase
	listPostsUsecase *post.ListPostsUsecase
}

func NewPostHandler(createUsecase *post.CreatePostUsecase, updateUsecase *post.UpdatePostUsecase, deleteUsecase *post.DeletePostUsecase, listPostsUsecase *post.ListPostsUsecase) *PostHandler {
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

	newPost, err := h.createUsecase.Execute(c.Request.Context(), post.CreatePostRequest{
		Title:    req.Title,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create post",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post": gin.H{
			"id":       newPost.ID,
			"title":    newPost.Title,
			"authorId": newPost.AuthorID,
			"content":  newPost.Content,
		},
	})
}

func (h *PostHandler) Update(c *gin.Context) {
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	updatedPost, err := h.updateUsecase.Execute(c.Request.Context(), post.UpdatePostRequest{
		ID:       req.ID,
		Title:    req.Title,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": gin.H{
			"id":       updatedPost.ID,
			"title":    updatedPost.Title,
			"authorId": updatedPost.AuthorID,
			"content":  updatedPost.Content,
		},
	})
}

func (h *PostHandler) Delete(c *gin.Context) {
	var req DeletePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	err := h.deleteUsecase.Execute(c.Request.Context(), post.DeletePostRequest{
		ID: req.ID,
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

	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"nextCursor": nextCursor,
	})
}
