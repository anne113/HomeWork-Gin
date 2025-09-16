package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	MySQLDSN    string
	JWTSecret   []byte
	TokenExpiry time.Duration
	ServerPort  string
)

func LoadConfig() {
	// 加载 .env 文件
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	MySQLDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	expiryStr := os.Getenv("TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "24" // 默认 24 小时
	}
	hours, err := strconv.Atoi(expiryStr)
	if err != nil {
		hours = 24
	}
	TokenExpiry = time.Duration(hours) * time.Hour

	ServerPort = os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = ":8080"
	}
}
