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

func TestPlaceOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockOrderService)
	orderController := controller.NewOrderController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.POST("/orders", orderController.PlaceOrder)

	// Case 1: Success
	mockService.On("PlaceOrder", uint(1), uint(10)).Return(nil)

	body := `{"address_id": 10}`
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Validation Error
	req2, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(`{}`))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Case 3: Service Error
	mockService.On("PlaceOrder", uint(1), uint(10)).Return(errors.New("failed"))
	req3, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(body))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
}

func TestGetOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockOrderService)
	orderController := controller.NewOrderController(mockService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
	})
	r.GET("/orders", orderController.GetOrders)

	mockService.On("GetOrders", uint(1)).Return([]model.Order{}, nil)

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Service Error
	mockService.On("GetOrders", uint(1)).Return(nil, errors.New("failed"))

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
