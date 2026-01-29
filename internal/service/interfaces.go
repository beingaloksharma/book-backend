package service

import "github.com/beingaloksharma/book-backend/internal/model"

type AuthServiceInterface interface {
	Signup(name, email, password string, role model.Role) error
	Login(email, password string) (string, error)
}

type UserServiceInterface interface {
	GetProfile(userID uint) (*model.User, error)
	AddAddress(userID uint, street, city, state, zip, country string) error
	GetAddresses(userID uint) ([]model.Address, error)
	GetAllUsers() ([]model.User, error)
}

type BookServiceInterface interface {
	CreateBook(title, author, description string, price float64, stock int) error
	UpdateBook(id uint, title, author, description string, price float64, stock int) error
	DeleteBook(id uint) error
	GetBook(id uint) (*model.Book, error)
	ListBooks() ([]model.Book, error)
}

type CartServiceInterface interface {
	AddToCart(userID, bookID uint, quantity int) error
	GetCart(userID uint) (*model.Cart, error)
}

type OrderServiceInterface interface {
	PlaceOrder(userID, addressID uint) error
	GetOrders(userID uint) ([]model.Order, error)
	GetAllOrders() ([]model.Order, error)
}
