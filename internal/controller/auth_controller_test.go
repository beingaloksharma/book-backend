package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/controller"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/service/mocks"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockAuthService)
	authController := controller.NewAuthController(mockService)

	r := gin.Default()
	r.POST("/signup", authController.Signup)

	// Case 1: Success
	mockService.On("Signup", "John", "john@example.com", "pass123", model.RoleUser).Return(nil).Once()

	body := `{"name":"John", "email":"john@example.com", "password":"pass123", "role":"USER"}`
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Bad Request
	body = `{"email":"john@example.com"}` // missing password and name
	req, _ = http.NewRequest("POST", "/signup", bytes.NewBufferString(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 3: Service Error
	mockService.On("Signup", "John", "john@example.com", "pass123", model.RoleUser).Return(errors.New("failed")).Once()
	body = `{"name":"John", "email":"john@example.com", "password":"pass123", "role":"USER"}`
	req, _ = http.NewRequest("POST", "/signup", bytes.NewBufferString(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockAuthService)
	authController := controller.NewAuthController(mockService)

	r := gin.Default()
	r.POST("/login", authController.Login)

	// Case 1: Success
	mockService.On("Login", "john@example.com", "pass123").Return("token123", nil).Once()

	body := `{"email":"john@example.com", "password":"pass123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token123")

	// Case 2: Unauthorized
	mockService.On("Login", "john@example.com", "wrong").Return("", errors.New("invalid")).Once()

	body = `{"email":"john@example.com", "password":"wrong"}`
	req, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
