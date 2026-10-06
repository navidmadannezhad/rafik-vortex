package utils

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func ResolveErrorMessage(msg string) gin.H {
	log.Printf("ERROR: %v", msg)
	return gin.H{
		"error": msg,
	}
}

func ResolveSuccessMessage(msg string) gin.H {
	return gin.H{
		"data": msg,
	}
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBased64From(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
