package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
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

	totalAmount := 0.0
	var orderItems []model.OrderItem

	for _, item := range cart.Items {
		price := item.Book.Price
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
	if err := s.OrderRepo.CreateOrder(order); err != nil {
		return err
	}

	// Clear Cart
	return s.CartRepo.ClearCart(cart.ID)
}

func (s *OrderService) GetOrders(userID uint) ([]model.Order, error) {
	return s.OrderRepo.FindByUserID(userID)
}
