package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// AuthMiddleware validates JWT tokens.
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", (*claims)["user_id"])
		c.Set("email", (*claims)["email"])
		c.Set("role", (*claims)["role"])

		c.Next()
	}
}

// RateLimitMiddleware limits request rate.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting
		c.Next()
	}
}

// LoggingMiddleware logs requests.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement request logging
		c.Next()
	}
}
