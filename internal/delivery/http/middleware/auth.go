package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userId"

func RequireUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-User-Id header"})
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(UserIDKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
