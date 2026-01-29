package model

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	gorm.Model
	UserID    uint        `json:"user_id"`
	AddressID uint        `json:"address_id"`
	Amount    float64     `json:"amount"`
	Status    OrderStatus `json:"status" gorm:"default:'PENDING'"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID  uint    `json:"order_id"`
	BookID   uint    `json:"book_id"`
	Book     Book    `json:"book"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // Captured price at time of order
}
