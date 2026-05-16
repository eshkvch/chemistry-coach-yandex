package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORS adds permissive CORS headers.
// Allowed origins are controlled via CORS_ALLOWED_ORIGINS env var (comma-separated).
// If not set, defaults to allowing all origins in development.
func CORS() gin.HandlerFunc {
	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowed == "" {
		allowed = "*"
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		// Determine whether to echo origin or use wildcard
		allowOrigin := allowed
		if allowed != "*" {
			// Check if origin is in the allowed list
			matched := false
			for _, o := range splitCSV(allowed) {
				if o == origin {
					matched = true
					break
				}
			}
			if matched {
				allowOrigin = origin
				c.Header("Vary", "Origin")
			} else {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func splitCSV(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			v := trim(s[start:i])
			if v != "" {
				result = append(result, v)
			}
			start = i + 1
		}
	}
	return result
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
