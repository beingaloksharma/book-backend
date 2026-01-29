package service

import (
	"errors"

	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/internal/repository"
)

type BookService struct {
	Repo *repository.BookRepository
}

func NewBookService() *BookService {
	return &BookService{Repo: repository.NewBookRepository()}
}

func (s *BookService) CreateBook(title, author, description string, price float64, stock int) error {
	book := &model.Book{
		Title:       title,
		Author:      author,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
	return s.Repo.CreateBook(book)
}

func (s *BookService) UpdateBook(id uint, title, author, description string, price float64, stock int) error {
	book, err := s.Repo.FindByID(id)
	if err != nil {
		return errors.New("book not found")
	}

	book.Title = title
	book.Author = author
	book.Description = description
	book.Price = price
	book.Stock = stock

	return s.Repo.UpdateBook(book)
}

func (s *BookService) DeleteBook(id uint) error {
	return s.Repo.DeleteBook(id)
}

func (s *BookService) GetBook(id uint) (*model.Book, error) {
	return s.Repo.FindByID(id)
}

func (s *BookService) ListBooks() ([]model.Book, error) {
	return s.Repo.FindAll()
}
