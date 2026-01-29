package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

func (r *BookRepository) CreateBook(book *model.Book) error {
	db := database.GetInstance()
	return db.Create(book).Error
}

func (r *BookRepository) UpdateBook(book *model.Book) error {
	db := database.GetInstance()
	return db.Save(book).Error
}

func (r *BookRepository) DeleteBook(id uint) error {
	db := database.GetInstance()
	return db.Delete(&model.Book{}, id).Error
}

func (r *BookRepository) FindByID(id uint) (*model.Book, error) {
	db := database.GetInstance()
	var book model.Book
	if err := db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindAll() ([]model.Book, error) {
	db := database.GetInstance()
	var books []model.Book
	if err := db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
