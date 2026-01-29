package repository

import "github.com/beingaloksharma/book-backend/internal/model"

type UserRepositoryInterface interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	AddAddress(address *model.Address) error
	GetAddresses(userID uint) ([]model.Address, error)
	FindAllUsers() ([]model.User, error)
}

type BookRepositoryInterface interface {
	CreateBook(book *model.Book) error
	UpdateBook(book *model.Book) error
	DeleteBook(id uint) error
	FindByID(id uint) (*model.Book, error)
	FindAll() ([]model.Book, error)
}

type CartRepositoryInterface interface {
	FindCartByUserID(userID uint) (*model.Cart, error)
	CreateCart(cart *model.Cart) error
	AddItem(item *model.CartItem) error
	UpdateItem(item *model.CartItem) error
	RemoveItem(itemID uint) error
	ClearCart(cartID uint) error
	FindItem(cartID, bookID uint) (*model.CartItem, error)
}

type OrderRepositoryInterface interface {
	CreateOrder(order *model.Order) error
	FindByUserID(userID uint) ([]model.Order, error)
	FindAllOrders() ([]model.Order, error)
	PlaceOrderTransaction(order *model.Order, cartItems []model.CartItem, cartID uint) error
}
