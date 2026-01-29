package repository

import (
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository() *BookRepository {
	return &BookRepository{DB: database.GetInstance()}
}

func (r *BookRepository) CreateBook(book *model.Book) error {
	return r.DB.Create(book).Error
}

func (r *BookRepository) UpdateBook(book *model.Book) error {
	return r.DB.Save(book).Error
}

func (r *BookRepository) DeleteBook(id uint) error {
	return r.DB.Delete(&model.Book{}, id).Error
}

func (r *BookRepository) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	if err := r.DB.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) FindAll() ([]model.Book, error) {
	var books []model.Book
	if err := r.DB.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
