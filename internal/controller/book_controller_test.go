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

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockBookService)
	bookController := controller.NewBookController(mockService)

	r := gin.Default()
	r.POST("/books", bookController.CreateBook)

	// Case 1: Success
	mockService.On("CreateBook", "Go", "Google", "Desc", 10.0, 5).Return(nil).Once()

	body := `{"title":"Go", "author":"Google", "description":"Desc", "price":10.0, "stock":5}`
	req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Validation Error
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString(`{}`))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Case 3: Service Error
	mockService.On("CreateBook", "Go", "Google", "Desc", 10.0, 5).Return(errors.New("failed")).Once()
	req, _ = http.NewRequest("POST", "/books", bytes.NewBufferString(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()

	mockService := new(mocks.MockBookService)
	bookController := controller.NewBookController(mockService)

	r := gin.Default()
	r.GET("/books/:id", bookController.GetBook)

	// Case 1: Success
	book := &model.Book{Title: "Go"}
	mockService.On("GetBook", uint(1)).Return(book, nil).Once()

	req, _ := http.NewRequest("GET", "/books/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Not Found
	mockService.On("GetBook", uint(2)).Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest("GET", "/books/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockBookService)
	bookController := controller.NewBookController(mockService)

	r := gin.Default()
	r.GET("/books", bookController.ListBooks)

	mockService.On("ListBooks").Return([]model.Book{{Title: "A"}}, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockBookService)
	bookController := controller.NewBookController(mockService)

	r := gin.Default()
	r.PUT("/books/:id", bookController.UpdateBook)

	// Update Success
	mockService.On("UpdateBook", uint(1), "Go", "Google", "Desc", 10.0, 5).Return(nil)

	body := `{"title":"Go", "author":"Google", "description":"Desc", "price":10.0, "stock":5}`
	req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Update Fail (validation)
	req2, _ := http.NewRequest("PUT", "/books/1", bytes.NewBufferString(`{}`))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// Update Fail (Service)
	mockService.On("UpdateBook", uint(1), "Go", "Google", "Desc", 10.0, 5).Return(errors.New("failed"))
	req3, _ := http.NewRequest("PUT", "/books/1", bytes.NewBufferString(body))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
}

func TestDeleteBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockBookService)
	bookController := controller.NewBookController(mockService)

	r := gin.Default()
	r.DELETE("/books/:id", bookController.DeleteBook)

	// Delete Success
	mockService.On("DeleteBook", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/books/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Delete Fail
	mockService.On("DeleteBook", uint(2)).Return(errors.New("failed"))
	req2, _ := http.NewRequest("DELETE", "/books/2", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
