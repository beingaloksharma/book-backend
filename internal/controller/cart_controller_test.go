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

func TestAddToCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockCartService)
	cartController := controller.NewCartController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.POST("/cart", cartController.AddToCart)

	// Case 1: Success
	mockService.On("AddToCart", uint(1), uint(10), 2).Return(nil)

	body := `{"book_id": 10, "quantity": 2}`
	req, _ := http.NewRequest("POST", "/cart", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Validation Error
	req2, _ := http.NewRequest("POST", "/cart", bytes.NewBufferString(`{}`))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Case 3: Service Error
	mockService.On("AddToCart", uint(1), uint(10), 2).Return(errors.New("failed"))
	req3, _ := http.NewRequest("POST", "/cart", bytes.NewBufferString(body))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
}

func TestGetCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockCartService)
	cartController := controller.NewCartController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.GET("/cart", cartController.GetCart)

	mockService.On("GetCart", uint(1)).Return(&model.Cart{}, nil)

	req, _ := http.NewRequest("GET", "/cart", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Cart not found
	mockService.On("GetCart", uint(2)).Return(nil, errors.New("not found"))
	r2 := gin.Default()
	r2.Use(func(c *gin.Context) {
		c.Set("user_id", uint(2))
	})
	r2.GET("/cart", cartController.GetCart)

	req2, _ := http.NewRequest("GET", "/cart", nil)
	w2 := httptest.NewRecorder()
	r2.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}
