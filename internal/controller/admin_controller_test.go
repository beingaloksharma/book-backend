package controller_test

import (
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

func TestAdminListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockUserService := new(mocks.MockUserService)
	mockOrderService := new(mocks.MockOrderService)
	adminController := controller.NewAdminController(mockUserService, mockOrderService)

	r := gin.Default()
	r.GET("/admin/users", adminController.ListUsers)

	mockUserService.On("GetAllUsers").Return([]model.User{}, nil)

	req, _ := http.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Service Error
	mockUserService.On("GetAllUsers").Return(nil, errors.New("failed"))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}

func TestAdminListOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUserService := new(mocks.MockUserService)
	mockOrderService := new(mocks.MockOrderService)
	adminController := controller.NewAdminController(mockUserService, mockOrderService)

	r := gin.Default()
	r.GET("/admin/orders", adminController.ListOrders)

	mockOrderService.On("GetAllOrders").Return([]model.Order{}, nil)

	req, _ := http.NewRequest("GET", "/admin/orders", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Service Error
	mockOrderService.On("GetAllOrders").Return(nil, errors.New("failed"))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
