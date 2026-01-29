package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
	CartRepo  *repository.CartRepository
	BookRepo  *repository.BookRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		OrderRepo: repository.NewOrderRepository(),
		CartRepo:  repository.NewCartRepository(),
		BookRepo:  repository.NewBookRepository(),
	}
}

func (s *OrderService) PlaceOrder(userID, addressID uint) error {
	// Get Cart
	cart, err := s.CartRepo.FindCartByUserID(userID)
	if err != nil || len(cart.Items) == 0 {
		return errors.New("cart is empty")
	}

	db := s.OrderRepo.DB

	// Start Transaction
	return db.Transaction(func(tx *gorm.DB) error {
		totalAmount := 0.0
		var orderItems []model.OrderItem

		for _, item := range cart.Items {
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

		order := &model.Order{
			UserID:    userID,
			AddressID: addressID,
			Amount:    totalAmount,
			Status:    model.OrderStatusPending,
			Items:     orderItems,
		}

		// Create Order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Clear Cart
		// Note: We use the transaction handler 'tx' to ensure atomic operations
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *OrderService) GetOrders(userID uint) ([]model.Order, error) {
	return s.OrderRepo.FindByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	return s.OrderRepo.FindAllOrders()
}
