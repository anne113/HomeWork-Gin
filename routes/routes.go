package routes

import (
	"HomeWork-Gin/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.POST("/posts", controllers.CreatePost)
	r.GET("/posts", controllers.GetPosts)

	return r
}
