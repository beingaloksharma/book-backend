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

func TestPlaceOrder(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockBookRepo)

	// Case 1: Cart Empty
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(&model.Cart{Items: []model.CartItem{}}, nil).Once()
	err := orderService.PlaceOrder(1, 1)
	assert.Error(t, err)
	assert.Equal(t, "cart is empty", err.Error())

	// Case 2: Success
	cart := &model.Cart{
		Model: gorm.Model{ID: 10},
		Items: []model.CartItem{
			{BookID: 1, Quantity: 2},
		},
	}
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(cart, nil).Once()
	mockOrderRepo.On("PlaceOrderTransaction", mock.AnythingOfType("*model.Order"), cart.Items, cart.ID).Return(nil).Once()

	err = orderService.PlaceOrder(1, 1)
	assert.NoError(t, err)

	mockOrderRepo.AssertExpectations(t)
	mockCartRepo.AssertExpectations(t)
}

func TestGetOrders(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockBookRepo)

	orders := []model.Order{{UserID: 1}}
	mockOrderRepo.On("FindByUserID", uint(1)).Return(orders, nil)

	result, err := orderService.GetOrders(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	mockOrderRepo.AssertExpectations(t)
}

func TestGetAllOrders(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockBookRepo)

	orders := []model.Order{{UserID: 1}}
	mockOrderRepo.On("FindAllOrders").Return(orders, nil)

	result, err := orderService.GetAllOrders()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))

	mockOrderRepo.AssertExpectations(t)
}

func TestPlaceOrder_RepoError(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockBookRepo)

	// Case 1: Cart Error
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(nil, errors.New("db error"))
	err := orderService.PlaceOrder(1, 1)
	assert.Error(t, err)

	// Case 2: Transaction Error
	cart := &model.Cart{
		Model: gorm.Model{ID: 10},
		Items: []model.CartItem{{BookID: 1, Quantity: 1}},
	}
	mockCartRepo.On("FindCartByUserID", uint(1)).Return(cart, nil)
	mockOrderRepo.On("PlaceOrderTransaction", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("tx error"))

	err = orderService.PlaceOrder(1, 1)
	assert.Error(t, err)
}

func TestGetOrders_Error(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockCartRepo := new(mocks.MockCartRepository)
	mockBookRepo := new(mocks.MockBookRepository)
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockBookRepo)

	mockOrderRepo.On("FindByUserID", uint(1)).Return(nil, errors.New("db error"))

	_, err := orderService.GetOrders(1)
	assert.Error(t, err)
}
