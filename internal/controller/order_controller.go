package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService *service.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{OrderService: service.NewOrderService()}
}

type PlaceOrderRequest struct {
	AddressID uint `json:"address_id" binding:"required"`
}

// PlaceOrder godoc
// @Summary Place an order
// @Description Place an order from the user's cart
// @Tags Order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PlaceOrderRequest true "Place Order Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/orders [post]
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req PlaceOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := c.OrderService.PlaceOrder(uid, req.AddressID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order placed successfully"})
}

// GetOrders godoc
// @Summary List user orders
// @Description Get a list of all orders for the logged-in user
// @Tags Order
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Order
// @Failure 500 {object} map[string]string
// @Router /api/orders [get]
func (c *OrderController) GetOrders(ctx *gin.Context) {
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

	orders, err := c.OrderService.GetOrders(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
