package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/middleware"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/beingaloksharma/book-backend/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	viper.Set("jwt.secret", "testsecret")
	token.Init()

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Case 1: No Header
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Case 2: Invalid Format
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Case 3: Invalid Token
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Case 4: Valid Token
	validToken, _ := token.GenerateToken(1, "USER")
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	viper.Set("jwt.secret", "testsecret")
	token.Init()

	r := gin.Default()
	// Mock Auth Middleware setting context
	r.Use(func(c *gin.Context) {
		tokenStr := c.GetHeader("X-Token")
		if tokenStr != "" {
			claims, _ := token.ValidateToken(tokenStr)
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
		}
		c.Next()
	})

	r.GET("/admin", middleware.RoleMiddleware("ADMIN"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Case 1: Forbidden (User trying to access Admin)
	userToken, _ := token.GenerateToken(1, "USER")
	req, _ := http.NewRequest("GET", "/admin", nil)
	req.Header.Set("X-Token", userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Case 2: Allowed (Admin accessing Admin)
	adminToken, _ := token.GenerateToken(2, "ADMIN")
	req, _ = http.NewRequest("GET", "/admin", nil)
	req.Header.Set("X-Token", adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
