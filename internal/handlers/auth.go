package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/models"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/services"
)

type AuthHandler struct {
    userService *services.UserService
    authService *services.AuthService
}

func NewAuthHandler(userService *services.UserService, authService *services.AuthService) *AuthHandler {
    return &AuthHandler{userService: userService, authService: authService}
}

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    hash, err := h.authService.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hash,
        Role:     models.RoleUser,
    }
    if err := h.userService.CreateUser(c.Request.Context(), user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.userService.GetUserByEmail(c.Request.Context(), req.Email)
    if err != nil || !h.authService.CheckPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }
    token, err := h.authService.GenerateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/register", h.Register)
    router.POST("/login", h.Login)
}