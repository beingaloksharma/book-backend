package mocks

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Signup(name, email, password string, role model.Role) error {
	args := m.Called(name, email, password, role)
	return args.Error(0)
}
func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

// MockUserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(userID uint) (*model.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserService) AddAddress(userID uint, street, city, state, zip, country string) error {
	args := m.Called(userID, street, city, state, zip, country)
	return args.Error(0)
}
func (m *MockUserService) GetAddresses(userID uint) ([]model.Address, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Address), args.Error(1)
}
func (m *MockUserService) GetAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

// MockBookService
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(title, author, description string, price float64, stock int) error {
	args := m.Called(title, author, description, price, stock)
	return args.Error(0)
}
func (m *MockBookService) UpdateBook(id uint, title, author, description string, price float64, stock int) error {
	args := m.Called(id, title, author, description, price, stock)
	return args.Error(0)
}
func (m *MockBookService) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockBookService) GetBook(id uint) (*model.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}
func (m *MockBookService) ListBooks() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

// MockCartService
type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) AddToCart(userID, bookID uint, quantity int) error {
	args := m.Called(userID, bookID, quantity)
	return args.Error(0)
}
func (m *MockCartService) GetCart(userID uint) (*model.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

// MockOrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) PlaceOrder(userID, addressID uint) error {
	args := m.Called(userID, addressID)
	return args.Error(0)
}
func (m *MockOrderService) GetOrders(userID uint) ([]model.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Order), args.Error(1)
}
func (m *MockOrderService) GetAllOrders() ([]model.Order, error) {
	args := m.Called()
	return args.Get(0).([]model.Order), args.Error(1)
}
