package controller

import (
	"net/http"

	"github.com/beingaloksharma/book-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{UserService: service.NewUserService()}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get profile of the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
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

	user, err := c.UserService.GetProfile(uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type AddressRequest struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	ZipCode string `json:"zip_code" binding:"required"`
	Country string `json:"country" binding:"required"`
}

// AddAddress godoc
// @Summary Add a new address
// @Description Add a new address for the user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddressRequest true "Address Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/addresses [post]
func (c *UserController) AddAddress(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	var req AddressRequest
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

	if err := c.UserService.AddAddress(uid, req.Street, req.City, req.State, req.ZipCode, req.Country); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Address added successfully"})
}

// GetAddresses godoc
// @Summary List user addresses
// @Description Get all addresses of the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Address
// @Failure 500 {object} map[string]string
// @Router /api/addresses [get]
func (c *UserController) GetAddresses(ctx *gin.Context) {
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

	addresses, err := c.UserService.GetAddresses(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, addresses)
}
