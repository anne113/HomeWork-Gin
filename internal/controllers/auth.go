package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"HomeWork-Gin/config"
	"HomeWork-Gin/internal/database"
	"HomeWork-Gin/internal/models"
	"HomeWork-Gin/internal/utils"
)

type RegisterPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var payload RegisterPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.JSON(c, http.StatusBadRequest, "invalid payload", nil)
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("password hashing failed")
		utils.JSON(c, http.StatusInternalServerError, "server error", nil)
		return
	}

	user := models.User{
		Username: payload.Username,
		Password: string(hashed),
		Email:    payload.Email,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		logrus.WithError(err).Error("create user failed")
		utils.JSON(c, http.StatusBadRequest, "could not create user (maybe duplicate username/email)", nil)
		return
	}
	// Do not expose password
	user.Password = ""
	utils.JSON(c, http.StatusCreated, "registered", user)
}

func Login(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.JSON(c, http.StatusBadRequest, "invalid payload", nil)
		return
	}
	var user models.User
	if err := database.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		utils.JSON(c, http.StatusUnauthorized, "invalid username or password", nil)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.JSON(c, http.StatusUnauthorized, "invalid username or password", nil)
		return
	}

	// create jwt
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(config.TokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(config.JWTSecret)
	if err != nil {
		logrus.WithError(err).Error("failed signing token")
		utils.JSON(c, http.StatusInternalServerError, "server error", nil)
		return
	}
	utils.JSON(c, http.StatusOK, "login successful", gin.H{"token": tokenStr})
}
