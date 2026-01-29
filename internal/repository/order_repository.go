package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
	"gorm.io/gorm"
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
