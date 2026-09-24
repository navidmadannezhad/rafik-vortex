package utils

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ResolveErrorMessage(msg string) gin.H {
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
