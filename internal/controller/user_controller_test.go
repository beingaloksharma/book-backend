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

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockUserService)
	userController := controller.NewUserController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.GET("/profile", userController.GetProfile)

	// Case 1: Success
	mockService.On("GetProfile", uint(1)).Return(&model.User{Name: "John"}, nil)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfile_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)
	userController := controller.NewUserController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.GET("/profile", userController.GetProfile)

	mockService.On("GetProfile", uint(1)).Return(nil, errors.New("failed"))

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockUserService)
	userController := controller.NewUserController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.POST("/addresses", userController.AddAddress)

	// Case 1: Success
	mockService.On("AddAddress", uint(1), "Street", "City", "State", "Zip", "Country").Return(nil).Once()

	body := `{"street":"Street", "city":"City", "state":"State", "zip_code":"Zip", "country":"Country"}`
	req, _ := http.NewRequest("POST", "/addresses", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Validation Error
	req2, _ := http.NewRequest("POST", "/addresses", bytes.NewBufferString(`{}`))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Case 3: Service Error
	mockService.On("AddAddress", uint(1), "Street", "City", "State", "Zip", "Country").Return(errors.New("failed")).Once()
	req3, _ := http.NewRequest("POST", "/addresses", bytes.NewBufferString(body))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
}

func TestGetAddresses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockUserService)
	userController := controller.NewUserController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.GET("/addresses", userController.GetAddresses)

	mockService.On("GetAddresses", uint(1)).Return([]model.Address{}, nil)

	req, _ := http.NewRequest("GET", "/addresses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
