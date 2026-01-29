package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{DB: database.GetInstance()}
}

func (r *CartRepository) FindCartByUserID(userID uint) (*model.Cart, error) {
	var cart model.Cart
	if err := r.DB.Where("user_id = ?", userID).Preload("Items.Book").First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) CreateCart(cart *model.Cart) error {
	return r.DB.Create(cart).Error
}

func (r *CartRepository) AddItem(item *model.CartItem) error {
	return r.DB.Create(item).Error
}

func (r *CartRepository) UpdateItem(item *model.CartItem) error {
	return r.DB.Save(item).Error
}

func (r *CartRepository) RemoveItem(itemID uint) error {
	return r.DB.Delete(&model.CartItem{}, itemID).Error
}

func (r *CartRepository) ClearCart(cartID uint) error {
	return r.DB.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

func (r *CartRepository) FindItem(cartID, bookID uint) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.DB.Where("cart_id = ? AND book_id = ?", cartID, bookID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
