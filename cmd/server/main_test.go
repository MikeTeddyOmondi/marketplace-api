package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/database"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/handlers"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/implementation"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Use in-memory SQLite for testing
	cfg := &config.DatabaseConfig{
		Driver: "sqlite",
		SQLite: config.SQLiteConfig{
			Path: ":memory:",
		},
	}

	dbManager, _ := database.NewManager(cfg)
	db, _ := dbManager.Connect()
	db.AutoMigrate(&models.User{}, &models.Product{})

	// Setup repositories and services
	userRepo := implementation.NewUserRepository(db)
	productRepo := implementation.NewProductRepository(db)

	constants := &config.Constants{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
		BusinessRules: config.BusinessRulesConfig{
			MaxProductsPerUser:   1000,
			DefaultProductStatus: "active",
		},
	}

	userService := services.NewUserService(userRepo, constants)
	productService := services.NewProductService(productRepo, userRepo, constants)

	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.New()
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	productHandler.RegisterRoutes(api)

	return router
}

func TestCreateProduct(t *testing.T) {
	router := setupTestRouter()

	// First create a user
	user := models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	userJSON, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Now create a product
	product := models.Product{
		Code:   "TEST001",
		Name:   "Test Product",
		Price:  100,
		UserID: 1,
	}
	productJSON, _ := json.Marshal(product)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "TEST001", response.Code)
	assert.Equal(t, "active", response.Status)
}
