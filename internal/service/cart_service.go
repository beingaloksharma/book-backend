package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"gorm.io/gorm"
)

type CartService struct {
	CartRepo *repository.CartRepository
	BookRepo *repository.BookRepository
}

func NewCartService() *CartService {
	return &CartService{
		CartRepo: repository.NewCartRepository(),
		BookRepo: repository.NewBookRepository(),
	}
}

func (s *CartService) AddToCart(userID, bookID uint, quantity int) error {
	// Check if book exists
	_, err := s.BookRepo.FindByID(bookID)
	if err != nil {
		return errors.New("book not found")
	}

	// Get or Create Cart
	cart, err := s.CartRepo.FindCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &model.Cart{UserID: userID}
			if err := s.CartRepo.CreateCart(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Check if item exists in cart
	item, err := s.CartRepo.FindItem(cart.ID, bookID)
	if err == nil {
		// Update quantity
		item.Quantity += quantity
		if item.Quantity <= 0 {
			return s.CartRepo.RemoveItem(item.ID)
		}
		return s.CartRepo.UpdateItem(item)
	}

	// Add new item
	newItem := &model.CartItem{
		CartID:   cart.ID,
		BookID:   bookID,
		Quantity: quantity,
	}
	return s.CartRepo.AddItem(newItem)
}

func (s *CartService) GetCart(userID uint) (*model.Cart, error) {
	return s.CartRepo.FindCartByUserID(userID)
}
