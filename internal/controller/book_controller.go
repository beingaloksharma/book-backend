package controller

import (
	"net/http"
	"strconv"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
)

type BookController struct {
	BookService service.BookServiceInterface
}

func NewBookController(bookService service.BookServiceInterface) *BookController {
	return &BookController{BookService: bookService}
}

type BookRequest struct {
	Title       string  `json:"title" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BookRequest true "Book Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/books [post]
func (c *BookController) CreateBook(ctx *gin.Context) {
	var req BookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := c.BookService.CreateBook(req.Title, req.Author, req.Description, req.Price, req.Stock); err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to create book")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Book created successfully"})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update an existing book (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param request body BookRequest true "Book Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/books/{id} [put]
func (c *BookController) UpdateBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid book ID")
		return
	}

	var req BookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := c.BookService.UpdateBook(uint(id), req.Title, req.Author, req.Description, req.Price, req.Stock); err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to update book")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/books/{id} [delete]
func (c *BookController) DeleteBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid book ID")
		return
	}

	if err := c.BookService.DeleteBook(uint(id)); err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to delete book")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Get details of a specific book
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Success 200 {object} model.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/books/{id} [get]
func (c *BookController) GetBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid book ID")
		return
	}

	book, err := c.BookService.GetBook(uint(id))
	if err != nil {
		logger.LogError(ctx, http.StatusNotFound, err, "Book not found")
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// ListBooks godoc
// @Summary List all books
// @Description Get a list of all available books
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Book
// @Failure 500 {object} map[string]string
// @Router /api/books [get]
func (c *BookController) ListBooks(ctx *gin.Context) {
	books, err := c.BookService.ListBooks()
	if err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to list books")
		return
	}

	ctx.JSON(http.StatusOK, books)
}
