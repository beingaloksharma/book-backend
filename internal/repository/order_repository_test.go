package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateOrder(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.OrderRepository{DB: db}

	order := &model.Order{UserID: 1, Amount: 100.0}

	mock.ExpectBegin()
	// Flexible query match
	mock.ExpectQuery(`INSERT INTO "orders"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateOrder(order)
	assert.NoError(t, err)
}

func TestFindOrdersByUserID(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.OrderRepository{DB: db}

	userID := uint(1)

	// Order Query
	orderRows := sqlmock.NewRows([]string{"id", "user_id", "amount", "status"}).
		AddRow(1, userID, 100.0, "PENDING")

	mock.ExpectQuery(`SELECT .* FROM "orders" WHERE .*user_id =`).
		WithArgs(userID).
		WillReturnRows(orderRows)

	// Preload Items
	mock.ExpectQuery(`SELECT .* FROM "order_items" WHERE "order_items"."order_id" =`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "book_id"}))

	orders, err := repo.FindByUserID(userID)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
}

func TestPlaceOrderTransaction(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.OrderRepository{DB: db}

	order := &model.Order{UserID: 1, Status: model.OrderStatusPending}
	cartItems := []model.CartItem{
		{
			Model:    gorm.Model{ID: 10},
			BookID:   100,
			Quantity: 2,
		},
	}
	cartID := uint(5)

	mock.ExpectBegin()

	// 1. Lock Book
	bookRows := sqlmock.NewRows([]string{"id", "title", "stock", "price"}).
		AddRow(100, "Go Book", 10, 50.0)
	mock.ExpectQuery(`SELECT .* FROM "books" WHERE .*"id" = .* FOR UPDATE`).
		WithArgs(100, 1).
		WillReturnRows(bookRows)

	// 2. Update Book Stock
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET`)).
		WillReturnResult(sqlmock.NewResult(100, 1))

	// 3. Create Order
	mock.ExpectQuery(`INSERT INTO "orders"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// 4. Create Order Items (One item)
	mock.ExpectQuery(`INSERT INTO "order_items"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// 5. Delete Cart Items (Soft Delete -> UPDATE)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "cart_items" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), cartID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := repo.PlaceOrderTransaction(order, cartItems, cartID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
