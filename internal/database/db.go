package database

import (
	"HomeWork-Gin/config"
	"HomeWork-Gin/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectAndMigrate() error {
	dsn := config.MySQLDSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Error("failed to connect database")
		return err
	}
	DB = db

	// 自动迁移
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		logrus.WithError(err).Error("failed to migrate database")
		return err
	}
	logrus.Info("database connected and migrated")
	return nil
}
