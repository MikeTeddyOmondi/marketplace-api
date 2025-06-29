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
	"github.com/MikeTeddyOmondi/marketplace-api/internal/middleware"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/implementation"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	cfg := &config.DatabaseConfig{
		Driver: "sqlite",
		SQLite: config.SQLiteConfig{
			Path: ":memory:",
		},
	}

	dbManager, _ := database.NewManager(cfg)
	db, _ := dbManager.Connect()
	db.AutoMigrate(&models.User{}, &models.Product{})

	constants := &config.Constants{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
		BusinessRules: config.BusinessRulesConfig{
			MaxProductsPerUser:   1000,
			DefaultProductStatus: "active",
		},
		Auth: config.AuthConfig{
			JWTSecret:       "test-secret-key-12345678901234567890123456789012",
			TokenExpiration: 72,
			PasswordCost:    4,
		},
	}

	authService := services.NewAuthService(constants)
	userRepo := implementation.NewUserRepository(db)
	productRepo := implementation.NewProductRepository(db)
	userService := services.NewUserService(userRepo, authService, constants)
	productService := services.NewProductService(productRepo, userRepo, constants)

	authHandler := handlers.NewAuthHandler(userService, authService)
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	userHandler.RegisterRoutes(protected, middleware.AuthMiddleware(authService), middleware.RoleMiddleware(models.RoleAdmin))
	productHandler.RegisterRoutes(protected)

	return router
}

func TestCreateProduct(t *testing.T) {
	router := setupTestRouter()

	// Register a user
	registerReq := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	registerJSON, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(registerJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Login to get JWT token
	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	loginJSON, _ := json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"]
	assert.NotEmpty(t, token)

	// Now create a product with the JWT token
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
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "TEST001", response.Code)
	assert.Equal(t, "active", response.Status)
}
