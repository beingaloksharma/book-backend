package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
)

type OrderService struct {
	OrderRepo repository.OrderRepositoryInterface
	CartRepo  repository.CartRepositoryInterface
	BookRepo  repository.BookRepositoryInterface
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, cartRepo repository.CartRepositoryInterface, bookRepo repository.BookRepositoryInterface) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
		CartRepo:  cartRepo,
		BookRepo:  bookRepo,
	}
}

func (s *OrderService) PlaceOrder(userID, addressID uint) error {
	// Get Cart
	cart, err := s.CartRepo.FindCartByUserID(userID)
	if err != nil || len(cart.Items) == 0 {
		return errors.New("cart is empty")
	}

	order := &model.Order{
		UserID:    userID,
		AddressID: addressID,
		Status:    model.OrderStatusPending,
	}

	// Use Transaction in Repository
	return s.OrderRepo.PlaceOrderTransaction(order, cart.Items, cart.ID)
}

func (s *OrderService) GetOrders(userID uint) ([]model.Order, error) {
	return s.OrderRepo.FindByUserID(userID)
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	return s.OrderRepo.FindAllOrders()
}
