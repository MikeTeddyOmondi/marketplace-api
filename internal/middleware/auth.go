package middleware

import (
    "slices"
	"net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/models"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/services"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            return
        }
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            return
        }
        token := tokenParts[1]
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)
        c.Next()
    }
}

func RoleMiddleware(allowedRoles ...models.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found"})
            return
        }
        role, ok := userRole.(models.Role)
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
            return
        }
        hasPermission := slices.Contains(allowedRoles, role)
        if !hasPermission {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            return
        }
        c.Next()
    }
}