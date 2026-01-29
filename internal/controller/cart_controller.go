package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	CartService *service.CartService
}

func NewCartController() *CartController {
	return &CartController{CartService: service.NewCartService()}
}

type AddToCartRequest struct {
	BookID   uint `json:"book_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cast float64 (from JWT claims usually) to uint
	// JWT claims are often float64 in generic maps
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	default:
		// Handle other cases or error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.CartService.AddToCart(uid, req.BookID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

func (c *CartController) GetCart(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	cart, err := c.CartService.GetCart(uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}
