package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
