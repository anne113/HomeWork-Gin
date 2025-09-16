package middleware

import (
	"HomeWork-Gin/config"
	"HomeWork-Gin/internal/database"
	"HomeWork-Gin/internal/models"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func extractTokenFromHeader(c *gin.Context) (string, error) {
	ah := c.GetHeader("Authorization")
	if ah == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.SplitN(ah, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": err.Error(), "data": nil})
			return
		}
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.JWTSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid token", "data": nil})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid token claims", "data": nil})
			return
		}
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid user id in token", "data": nil})
			return
		}
		userID := uint(userIDFloat)

		// load user from DB optionally to ensure exists
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "user not found", "data": nil})
			return
		}
		// put user info in context
		c.Set("currentUser", user)
		c.Next()
	}
}
