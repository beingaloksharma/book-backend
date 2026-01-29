package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

func TestFindUserByEmail(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	email := "john@example.com"

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "role"}).
		AddRow(1, time.Now(), time.Now(), nil, "John", email, "USER")

	// Loose regex: match SELECT ... FROM "users" ... email =
	mock.ExpectQuery(`SELECT .* FROM "users" WHERE .*email =`).
		WithArgs(email, 1). // Email and Limit
		WillReturnRows(rows)

	user, err := repo.FindByEmail(email)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindUserByID(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	id := uint(1)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email", "role"}).
		AddRow(id, time.Now(), time.Now(), nil, "John", "john@example.com", "USER")

	mock.ExpectQuery(`SELECT .* FROM "users" WHERE .*"id" =`).
		WithArgs(id, 1). // ID and Limit
		WillReturnRows(rows)

	user, err := repo.FindByID(id)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, id, user.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	user := &model.User{Name: "John"}

	mock.ExpectBegin()
	// Flexible
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
}

func TestAddAddress(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	addr := &model.Address{UserID: 1, City: "City"}

	// Flexible
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "addresses"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.AddAddress(addr)
	assert.NoError(t, err)
}

func TestGetAddresses(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "city"}).AddRow(1, 1, "City")

	mock.ExpectQuery(`SELECT .* FROM "addresses" WHERE .*user_id =`).
		WithArgs(1).
		WillReturnRows(rows)

	addrs, err := repo.GetAddresses(1)
	require.NoError(t, err)
	assert.Len(t, addrs, 1)
}

func TestFindAllUsers(t *testing.T) {
	db, mock := NewMockDB()
	repo := &repository.UserRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John")

	mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(rows)

	users, err := repo.FindAllUsers()
	require.NoError(t, err)
	assert.Len(t, users, 1)
}
