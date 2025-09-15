package controllers

import (
	"blog/config"
	"blog/models"
	"blog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := utils.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效token"})
		return
	}

	var input struct {
		Title   string
		Content string
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{Title: input.Title, Content: input.Content, UserID: claims.UserID}
	config.DB.Create(&post)
	c.JSON(http.StatusOK, gin.H{"message": "文章创建成功", "post": post})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("User").Preload("Comments").Find(&posts)
	c.JSON(http.StatusOK, posts)
}
