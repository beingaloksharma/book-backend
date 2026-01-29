package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	CartService service.CartServiceInterface
}

func NewCartController(cartService service.CartServiceInterface) *CartController {
	return &CartController{CartService: cartService}
}

type AddToCartRequest struct {
	BookID   uint `json:"book_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a book to user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddToCartRequest true "Add To Cart Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart [post]
func (c *CartController) AddToCart(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid request body")
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
		logger.LogError(ctx, http.StatusInternalServerError, nil, "Invalid user ID")
		return
	}

	if err := c.CartService.AddToCart(uid, req.BookID, req.Quantity); err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to add item to cart")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

// GetCart godoc
// @Summary Get user cart
// @Description Get current user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Cart
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/cart [get]
func (c *CartController) GetCart(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	default:
		logger.LogError(ctx, http.StatusInternalServerError, nil, "Invalid user ID")
		return
	}

	cart, err := c.CartService.GetCart(uid)
	if err != nil {
		logger.LogError(ctx, http.StatusNotFound, err, "Cart not found")
		return
	}

	ctx.JSON(http.StatusOK, cart)
}
