package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id" gorm:"unique"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	gorm.Model
	CartID   uint `json:"cart_id"`
	BookID   uint `json:"book_id"`
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}
