package controllers

import (
	"HomeWork-Gin/internal/database"
	"HomeWork-Gin/internal/models"
	"HomeWork-Gin/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostPayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreatePost(c *gin.Context) {
	var payload PostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.JSON(c, http.StatusBadRequest, "invalid payload", nil)
		return
	}
	current := c.MustGet("currentUser").(models.User)
	post := models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  current.ID,
	}
	if err := database.DB.Create(&post).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed create post", nil)
		return
	}
	utils.JSON(c, http.StatusCreated, "post created", post)
}

func ListPosts(c *gin.Context) {
	var posts []models.Post
	if err := database.DB.Preload("User").Find(&posts).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed to fetch posts", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "ok", posts)
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Preload("User").Preload("Comments").First(&post, id).Error; err != nil {
		utils.JSON(c, http.StatusNotFound, "post not found", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "ok", post)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		utils.JSON(c, http.StatusNotFound, "post not found", nil)
		return
	}
	current := c.MustGet("currentUser").(models.User)
	if post.UserID != current.ID {
		utils.JSON(c, http.StatusForbidden, "not allowed", nil)
		return
	}
	var payload PostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.JSON(c, http.StatusBadRequest, "invalid payload", nil)
		return
	}
	post.Title = payload.Title
	post.Content = payload.Content
	if err := database.DB.Save(&post).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed to update", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "updated", post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		utils.JSON(c, http.StatusNotFound, "post not found", nil)
		return
	}
	current := c.MustGet("currentUser").(models.User)
	if post.UserID != current.ID {
		utils.JSON(c, http.StatusForbidden, "not allowed", nil)
		return
	}
	if err := database.DB.Delete(&post).Error; err != nil {
		utils.JSON(c, http.StatusInternalServerError, "failed to delete", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "deleted", nil)
}
