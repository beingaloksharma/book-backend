package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logrus.Infof("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown: ", err)
	}

	logrus.Info("Server exiting")
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
