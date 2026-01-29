package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	OrderService service.OrderServiceInterface
}

func NewOrderController(orderService service.OrderServiceInterface) *OrderController {
	return &OrderController{OrderService: orderService}
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
		logger.LogError(ctx, http.StatusBadRequest, err, "Invalid request body")
		return
	}

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

	if err := c.OrderService.PlaceOrder(uid, req.AddressID); err != nil {
		logger.LogError(ctx, http.StatusBadRequest, err, "Failed to place order")
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
		logger.LogError(ctx, http.StatusInternalServerError, nil, "Invalid user ID")
		return
	}

	orders, err := c.OrderService.GetOrders(uid)
	if err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to fetch orders")
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
