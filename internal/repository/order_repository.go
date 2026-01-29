package repository

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{DB: database.GetInstance()}
}

func (r *OrderRepository) CreateOrder(order *model.Order) error {
	db := database.GetInstance()
	return db.Create(order).Error
}

func (r *OrderRepository) FindByUserID(userID uint) ([]model.Order, error) {
	db := database.GetInstance()
	var orders []model.Order
	if err := db.Where("user_id = ?", userID).Preload("Items.Book").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindAllOrders() ([]model.Order, error) {
	db := database.GetInstance()
	var orders []model.Order
	if err := db.Preload("Items.Book").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) PlaceOrderTransaction(order *model.Order, cartItems []model.CartItem, cartID uint) error {
	// Start Transaction
	return r.DB.Transaction(func(tx *gorm.DB) error {
		totalAmount := 0.0
		var orderItems []model.OrderItem

		for _, item := range cartItems {
			// Lock book row for update to prevent race conditions
			var book model.Book
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, item.BookID).Error; err != nil {
				return err
			}

			if book.Stock < item.Quantity {
				return errors.New("insufficient stock for book: " + book.Title)
			}

			// Deduct Stock
			book.Stock -= item.Quantity
			if err := tx.Save(&book).Error; err != nil {
				return err
			}

			price := book.Price
			totalAmount += price * float64(item.Quantity)
			orderItems = append(orderItems, model.OrderItem{
				BookID:   item.BookID,
				Quantity: item.Quantity,
				Price:    price,
			})
		}

		order.Amount = totalAmount
		order.Items = orderItems

		// Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Clear Cart
		if err := tx.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}
