package middleware

import (
	"net/http"
	"strings"

	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/beingaloksharma/book-backend/utils/token"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.LogError(c, http.StatusUnauthorized, nil, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.LogError(c, http.StatusUnauthorized, nil, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := token.ValidateToken(parts[1])
		if err != nil {
			logger.LogError(c, http.StatusUnauthorized, err, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			logger.LogError(c, http.StatusForbidden, nil, "Forbidden: insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
