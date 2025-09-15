package main

import (
	"HomeWork-Gin/config"
	"HomeWork-Gin/models"
	"HomeWork-Gin/routes"
)

func main() {
	config.InitDB()
	// 自动建表
	config.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	r := routes.SetupRouter()
	r.Run(":8080")
}
