package service_test

import (
	"errors"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository/mocks"
	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBook(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	mockRepo.On("CreateBook", mock.AnythingOfType("*model.Book")).Return(nil)

	err := bookService.CreateBook("Go", "Google", "Lang", 10.0, 5)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetBook(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	// Success
	book := &model.Book{Title: "Go"}
	mockRepo.On("FindByID", uint(1)).Return(book, nil)

	result, err := bookService.GetBook(1)
	assert.NoError(t, err)
	assert.Equal(t, "Go", result.Title)

	// Error
	mockRepo.On("FindByID", uint(2)).Return(nil, errors.New("not found"))

	_, err = bookService.GetBook(2)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestListBooks(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	bookService := service.NewBookService(mockRepo)

	books := []model.Book{{Title: "A"}, {Title: "B"}}
	mockRepo.On("FindAll").Return(books, nil)

	result, err := bookService.ListBooks()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	mockRepo.AssertExpectations(t)
}
