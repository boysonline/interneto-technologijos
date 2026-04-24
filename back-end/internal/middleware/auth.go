package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedKey := c.GetHeader("X-API-KEY")
		requiredKey := os.Getenv("API_KEY")

		if receivedKey == "" || receivedKey != requiredKey {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Invalid API Key"})
			return
		}

		c.Next()
	}
}
