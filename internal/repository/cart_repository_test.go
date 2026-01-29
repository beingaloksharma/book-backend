package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCartByUserID(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.CartRepository{DB: db}

	userID := uint(1)

	// Expect Cart Query
	cartRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id"}).
		AddRow(1, time.Now(), time.Now(), nil, userID)

	// Expect Items Preload (This is tricky with GORM + SQLMock as GORM does separate query)
	// Usually: SELECT * FROM carts WHERE user_id = ?
	// Then: SELECT * FROM cart_items WHERE cart_id IN (1)

	mock.ExpectQuery(`SELECT .* FROM "carts" WHERE .*user_id =`).
		WithArgs(userID, 1).
		WillReturnRows(cartRows)

	// Mock Items query if preload happens
	// Assuming no items for simplicity to pass basic "Found" check
	// If Preload is strictly required, we need another ExpectQuery
	mock.ExpectQuery(`SELECT .* FROM "cart_items" WHERE "cart_items"."cart_id" =`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id"}))

	cart, err := repo.FindCartByUserID(userID)
	require.NoError(t, err)
	assert.Equal(t, userID, cart.UserID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCart(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.CartRepository{DB: db}

	cart := &model.Cart{UserID: 1}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "carts"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), cart.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateCart(cart)
	assert.NoError(t, err)
}

func TestAddItem(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.CartRepository{DB: db}

	item := &model.CartItem{CartID: 1, BookID: 2, Quantity: 1}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "cart_items"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), item.CartID, item.BookID, item.Quantity).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.AddItem(item)
	assert.NoError(t, err)
}

func TestFindItem(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.CartRepository{DB: db}

	cartID := uint(1)
	bookID := uint(2)

	rows := sqlmock.NewRows([]string{"id", "cart_id", "book_id", "quantity"}).
		AddRow(1, cartID, bookID, 5)

	mock.ExpectQuery(`SELECT .* FROM "cart_items" WHERE .*cart_id = .* AND book_id =`).
		WithArgs(cartID, bookID, 1).
		WillReturnRows(rows)

	item, err := repo.FindItem(cartID, bookID)
	require.NoError(t, err)
	assert.Equal(t, 5, item.Quantity)
}
