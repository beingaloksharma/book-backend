package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (r *CartRepository) FindCartByUserID(userID uint) (*model.Cart, error) {
	db := database.GetInstance()
	var cart model.Cart
	if err := db.Where("user_id = ?", userID).Preload("Items.Book").First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(cart *model.Cart) error {
	db := database.GetInstance()
	return db.Create(cart).Error
}

func (r *CartRepository) AddItem(item *model.CartItem) error {
	db := database.GetInstance()
	return db.Create(item).Error
}

func (r *CartRepository) UpdateItem(item *model.CartItem) error {
	db := database.GetInstance()
	return db.Save(item).Error
}

func (r *CartRepository) RemoveItem(itemID uint) error {
	db := database.GetInstance()
	return db.Delete(&model.CartItem{}, itemID).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	db := database.GetInstance()
	return db.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

func (r *CartRepository) FindItem(cartID, bookID uint) (*model.CartItem, error) {
	db := database.GetInstance()
	var item model.CartItem
	if err := db.Where("cart_id = ? AND book_id = ?", cartID, bookID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
