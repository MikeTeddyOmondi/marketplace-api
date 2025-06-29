package services

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/config"
    "github.com/MikeTeddyOmondi/marketplace-api/internal/models"
)

type AuthService struct {
    jwtSecret    string
    tokenExpiry  time.Duration
    passwordCost int
}

func NewAuthService(constants *config.Constants) *AuthService {
    return &AuthService{
        jwtSecret:    constants.Auth.JWTSecret,
        tokenExpiry:  time.Duration(constants.Auth.TokenExpiration) * time.Hour,
        passwordCost: constants.Auth.PasswordCost,
    }
}

type Claims struct {
    UserID uint        `json:"user_id"`
    Role   models.Role `json:"role"`
    jwt.RegisteredClaims
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
    expirationTime := time.Now().Add(s.tokenExpiry)
    claims := &Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            Subject:   user.Email,
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}

func (s *AuthService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.passwordCost)
    return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
