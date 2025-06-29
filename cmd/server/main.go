package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/database"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/handlers"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/implementation"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.App.Server.Mode)

	// Initialize database
	dbManager, err := database.NewManager(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database manager: %v", err)
	}

	db, err := dbManager.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := implementation.NewUserRepository(db)
	productRepo := implementation.NewProductRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, &cfg.Constants)
	productService := services.NewProductService(productRepo, userRepo, &cfg.Constants)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	// Setup router
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.App.CORS.AllowedOrigins
	corsConfig.AllowMethods = cfg.App.CORS.AllowedMethods
	corsConfig.AllowHeaders = cfg.App.CORS.AllowedHeaders
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"database":  "connected",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	productHandler.RegisterRoutes(api)

	// Start server
	srv := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.App.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.App.Server.Timeout) * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server starting on %s", cfg.GetAddress())

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close database connection
	if err := dbManager.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited")
}
