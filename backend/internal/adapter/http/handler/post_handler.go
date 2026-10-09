package handler

import (
	"blog-aws-backend/internal/usecase/post"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageKey string `json:"imageKey" binding:"required"`
}

type UpdatePostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageKey string `json:"imageKey" binding:"required"`
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

	userID := c.GetString("userID")
	newPost, err := h.createUsecase.Execute(c.Request.Context(), post.CreatePostRequest{
		Title:    req.Title,
		AuthorID: userID,
		Content:  req.Content,
		ImageKey: req.ImageKey,
	})
	fmt.Println("error ================================", err)
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
			"imageKey": newPost.ImageKey,
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

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	userID := c.GetString("userID")
	updatedPost, err := h.updateUsecase.Execute(c.Request.Context(), post.UpdatePostRequest{
		ID:       id,
		Title:    req.Title,
		AuthorID: userID,
		Content:  req.Content,
		ImageKey: req.ImageKey,
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
			"imageKey": updatedPost.ImageKey,
		},
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

	c.JSON(http.StatusOK, gin.H{
		"posts":      posts,
		"nextCursor": nextCursor,
	})
}
