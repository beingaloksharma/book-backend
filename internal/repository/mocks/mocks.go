package mocks

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) AddAddress(address *model.Address) error {
	args := m.Called(address)
	return args.Error(0)
}
func (m *MockUserRepository) GetAddresses(userID uint) ([]model.Address, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Address), args.Error(1)
}
func (m *MockUserRepository) FindAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

// MockBookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}
func (m *MockBookRepository) UpdateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}
func (m *MockBookRepository) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockBookRepository) FindByID(id uint) (*model.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}
func (m *MockBookRepository) FindAll() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

// MockCartRepository
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FindCartByUserID(userID uint) (*model.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}
func (m *MockCartRepository) CreateCart(cart *model.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}
func (m *MockCartRepository) AddItem(item *model.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}
func (m *MockCartRepository) UpdateItem(item *model.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}
func (m *MockCartRepository) RemoveItem(itemID uint) error {
	args := m.Called(itemID)
	return args.Error(0)
}
func (m *MockCartRepository) ClearCart(cartID uint) error {
	args := m.Called(cartID)
	return args.Error(0)
}
func (m *MockCartRepository) FindItem(cartID, bookID uint) (*model.CartItem, error) {
	args := m.Called(cartID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CartItem), args.Error(1)
}

// MockOrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}
func (m *MockOrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Order), args.Error(1)
}
func (m *MockOrderRepository) FindAllOrders() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}
func (m *MockOrderRepository) PlaceOrderTransaction(order *model.Order, cartItems []model.CartItem, cartID uint) error {
	args := m.Called(order, cartItems, cartID)
	return args.Error(0)
}
