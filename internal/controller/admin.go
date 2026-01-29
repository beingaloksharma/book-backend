package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	UserService  *service.UserService
	OrderService *service.OrderService
}

func NewAdminController() *AdminController {
	return &AdminController{
		UserService:  service.NewUserService(),
		OrderService: service.NewOrderService(),
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all registered users (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.User
// @Failure 500 {object} map[string]string
// @Router /api/admin/users [get]
func (c *AdminController) ListUsers(ctx *gin.Context) {
	users, err := c.UserService.GetAllUsers()
	if err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to fetch users")
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// ListOrders godoc
// @Summary List all orders
// @Description Get a list of all orders (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Order
// @Failure 500 {object} map[string]string
// @Router /api/admin/orders [get]
func (c *AdminController) ListOrders(ctx *gin.Context) {
	orders, err := c.OrderService.GetAllOrders()
	if err != nil {
		logger.LogError(ctx, http.StatusInternalServerError, err, "Failed to fetch orders")
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
