package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindBookByID(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.BookRepository{DB: db}

	id := uint(1)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "author", "price", "stock"}).
		AddRow(id, time.Now(), time.Now(), nil, "Go", "Google", 10.0, 5)

	mock.ExpectQuery(`SELECT .* FROM "books" WHERE .*"id" =`).
		WithArgs(id, 1). // ID and Limit
		WillReturnRows(rows)

	book, err := repo.FindByID(id)
	require.NoError(t, err)
	require.NotNil(t, book)
	assert.Equal(t, "Go", book.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllBooks(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.BookRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title"}).
		AddRow(1, time.Now(), time.Now(), nil, "Book A").
		AddRow(2, time.Now(), time.Now(), nil, "Book B")

	mock.ExpectQuery(`SELECT .* FROM "books"`).
		WillReturnRows(rows)

	books, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, books, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateBook(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.BookRepository{DB: db}

	book := &model.Book{Title: "New Book", Price: 20.0}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "books"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateBook(book)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.BookRepository{DB: db}

	book := &model.Book{Title: "Updated Book"}
	// GORM's Save updates all fields. We'll simplify the expectation for now or assume it updates the row.
	// Since Save can be an INSERT or UPDATE depending on ID presence, let's assume valid ID > 0
	book.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateBook(book)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBook(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.BookRepository{DB: db}

	id := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "books" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteBook(id)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
