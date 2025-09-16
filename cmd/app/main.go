package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"HomeWork-Gin/config"
	"HomeWork-Gin/internal/controllers"
	"HomeWork-Gin/internal/database"
	"HomeWork-Gin/internal/middleware"
)

func main() {
	// 日志配置
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// 加载配置
	config.LoadConfig()

	// 连接数据库 + 自动迁移
	if err := database.ConnectAndMigrate(); err != nil {
		logrus.Fatal("failed to connect to db: ", err)
	}

	r := gin.Default()

	// Auth
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// 需要认证的 API
	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		auth.POST("/posts", controllers.CreatePost)
		auth.PUT("/posts/:id", controllers.UpdatePost)
		auth.DELETE("/posts/:id", controllers.DeletePost)

		auth.POST("/posts/:post_id/comments", controllers.CreateComment)
	}

	// 公共 API
	r.GET("/api/posts", controllers.ListPosts)
	r.GET("/api/posts/:id", controllers.GetPost)
	r.GET("/api/posts/:post_id/comments", controllers.ListComments)

	// 启动服务
	logrus.Infof("starting server on %s", config.ServerPort)
	if err := r.Run(config.ServerPort); err != nil {
		logrus.Fatal(err)
	}
}
