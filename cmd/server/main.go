package main

import (
	"flag"
	"fmt"

	_ "github.com/beingaloksharma/book-backend/docs"
	"github.com/beingaloksharma/book-backend/internal/controller"
	"github.com/beingaloksharma/book-backend/internal/middleware"
	"github.com/beingaloksharma/book-backend/internal/model"
	"github.com/beingaloksharma/book-backend/utils/database"
	"github.com/beingaloksharma/book-backend/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Book Store API
// @version 1.0
// @description REST API for a Book Store application managing users, books, carts, and orders.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Logger - utilizing the custom wrapper
	logger.GinLogger()
	// 1. Define Flag
	// Improved description for the help text
	configFilePath := flag.String("config-path", "config/", "Path to the configuration directory")
	// 2. Parse Flags
	flag.Parse()
	// 3. Load Configuration
	// Pass the value of the pointer (*configFilePath) directly to the function
	loadConfig(*configFilePath)
	r := gin.Default()
	setupDatabase(r)

	// Init Controllers
	authController := controller.NewAuthController()

	// Routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Book Routes
	bookController := controller.NewBookController()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// User, Cart, Order Controllers
	userController := controller.NewUserController()
	cartController := controller.NewCartController()
	orderController := controller.NewOrderController()

	// Public/User Book Routes
	api.GET("/books", bookController.ListBooks)
	api.GET("/books/:id", bookController.GetBook)

	// User Routes
	api.GET("/profile", userController.GetProfile)
	api.POST("/addresses", userController.AddAddress)
	api.GET("/addresses", userController.GetAddresses)

	// Cart Routes
	api.POST("/cart", cartController.AddToCart)
	api.GET("/cart", cartController.GetCart) // Review cart

	// Order Routes
	api.POST("/orders", orderController.PlaceOrder) // Make order
	api.GET("/orders", orderController.GetOrders)

	// Admin Book Routes
	admin := api.Group("/admin")
	admin.Use(middleware.RoleMiddleware("ADMIN"))
	{
		admin.POST("/books", bookController.CreateBook)
		admin.PUT("/books/:id", bookController.UpdateBook)
		admin.DELETE("/books/:id", bookController.DeleteBook)
		admin.GET("/profile", userController.GetProfile) // reusing user profile for admin
	}

	fmt.Println("Hello World")

	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}
	logrus.Infof("Starting server on port %s", port)
	if err := r.Run(port); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig - Load the config parameters
// Accepts path as an argument to avoid reliance on global variables
func loadConfig(path string) {
	viper.SetConfigName("app-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)           // Default config path passed
	viper.AddConfigPath(".")            // Current directory
	viper.AddConfigPath("../../config") // Parent config directory (if running from cmd/server)
	viper.AddConfigPath("./config")     // Config directory in current

	if err := viper.ReadInConfig(); err != nil {
		// Type assertion to check specifically for FileNotFound
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Using the global 'log' wrapper for consistency
			logrus.Errorf("No configuration file found at %s. Using defaults if available.", path)
		} else {
			// FIXED: Use 'err' here, not 'readErr'
			// Using 'log' instead of 'logrus' directly
			logrus.Errorf("Error reading config file: %s", err)
		}
	} else {
		logrus.Infof("Configuration loaded successfully from %s", path)
		logrus.Info("Configuration Key and Value are printed below")
		logrus.Info("-------------------------")
		for key, val := range viper.AllSettings() {
			logrus.Infof("%s: %v", key, val)
		}

	}
}

// Database Connection
func setupDatabase(r *gin.Engine) {
	database.GetInstance()
	database.Migrate(
		&model.User{},
		&model.Book{},
		&model.Address{},
		&model.Cart{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	)
}
