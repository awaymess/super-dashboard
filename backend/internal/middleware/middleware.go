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

// OptionalAuthMiddleware validates JWT tokens but doesn't require them.
// Useful for endpoints that work differently for authenticated vs anonymous users.
func OptionalAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Set user info in context if token is valid
		c.Set("user_id", (*claims)["user_id"])
		c.Set("email", (*claims)["email"])
		c.Set("role", (*claims)["role"])

		c.Next()
	}
}

// RoleMiddleware requires a specific role to access the endpoint.
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, r := range requiredRoles {
			if role == r {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware requires admin role to access the endpoint.
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// LoggingMiddleware logs requests.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement request logging
		c.Next()
	}
}
