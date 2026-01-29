package service_test

import (
	"errors"
	"testing"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository/mocks"
	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddToCart(t *testing.T) {
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	cartService := service.NewCartService(mockCartRepo, mockBookRepo)

	// Case 1: Book Not Found
	mockBookRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found")).Once()
	err := cartService.AddToCart(1, 1, 1)
	assert.Error(t, err)

	// Case 2: Cart not found, create new cart
	mockBookRepo.On("FindByID", uint(1)).Return(&model.Book{}, nil)
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()
	mockCartRepo.On("CreateCart", mock.AnythingOfType("*model.Cart")).Run(func(args mock.Arguments) {
		cart := args.Get(0).(*model.Cart)
		cart.ID = 10 // simulate DB ID
	}).Return(nil).Once()

	// Then item check fails (new cart), so add item
	mockCartRepo.On("FindItem", uint(10), uint(1)).Return(nil, errors.New("not found")).Once()
	mockCartRepo.On("AddItem", mock.AnythingOfType("*model.CartItem")).Return(nil).Once()

	err = cartService.AddToCart(1, 1, 1)
	assert.NoError(t, err)

	// Case 3: Cart exists, item exists, update quantity
	mockBookRepo.On("FindByID", uint(1)).Return(&model.Book{}, nil)
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(&model.Cart{
		Model:  gorm.Model{ID: 10},
		UserID: 1,
	}, nil).Once()
	mockCartRepo.On("FindItem", uint(10), uint(1)).Return(&model.CartItem{
		Model:    gorm.Model{ID: 5},
		Quantity: 1,
	}, nil).Once()
	mockCartRepo.On("UpdateItem", mock.AnythingOfType("*model.CartItem")).Return(nil).Once()

	err = cartService.AddToCart(1, 1, 1) // +1 quantity
	assert.NoError(t, err)

	mockCartRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestGetCart(t *testing.T) {
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	cartService := service.NewCartService(mockCartRepo, mockBookRepo)

	cart := &model.Cart{UserID: 1}
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(cart, nil)

	result, err := cartService.GetCart(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)

	mockCartRepo.AssertExpectations(t)
}

func TestGetCart_Error(t *testing.T) {
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	cartService := service.NewCartService(mockCartRepo, mockBookRepo)

	mockCartRepo.On("FindCartByUserID", uint(1)).Return(nil, errors.New("db error"))

	_, err := cartService.GetCart(1)
	assert.Error(t, err)
}

func TestAddToCart_RepoError(t *testing.T) {
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	cartService := service.NewCartService(mockCartRepo, mockBookRepo)

	mockBookRepo.On("FindByID", uint(1)).Return(&model.Book{}, nil)
	// Fail finding cart
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(nil, errors.New("db error")) // Not ErrRecordNotFound

	err := cartService.AddToCart(1, 1, 1)
	assert.Error(t, err)
}
