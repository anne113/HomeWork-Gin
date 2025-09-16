package controllers

import (
	"net/http"

	"HomeWork-Gin/internal/database"
	"HomeWork-Gin/internal/models"
	"HomeWork-Gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type CommentPayload struct {
	Content string `json:"content" binding:"required"`
}

func CreateComment(c *gin.Context) {
	postID := c.Param("post_id")
	var payload CommentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.JSON(c, http.StatusBadRequest, "invalid payload", nil)
		return
	}
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.JSON(c, http.StatusNotFound, "post not found", nil)
		return
	}
	current := c.MustGet("currentUser").(models.User)
	comment := models.Comment{
		Content: payload.Content,
		UserID:  current.ID,
		PostID:  post.ID,
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed create comment", nil)
		return
	}
	utils.JSON(c, http.StatusCreated, "comment created", comment)
}

func ListComments(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []models.Comment
	if err := database.DB.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed to fetch comments", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "ok", comments)
}
