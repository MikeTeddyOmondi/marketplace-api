# Modular Repository Pattern with GORM, Gin & YAML

A comprehensive guide to implementing a modular repository pattern using GORM ORM, Gin web framework, and YAML configuration. This architecture allows for easy database switching and maintains clean separation of concerns.

## Table of Contents

1. [Project Structure](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#project-structure)
2. [Dependencies](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#dependencies)
3. [Configuration](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#configuration)
4. [Database Interface](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#database-interface)
5. [Models](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#models)
6. [Repository Layer](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#repository-layer)
7. [Service Layer](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#service-layer)
8. [API Handlers](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#api-handlers)
9. [Main Application](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#main-application)
10. [Usage Examples](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#usage-examples)
11. [Database Migration](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#database-migration)
12. [Testing](https://claude.ai/chat/2d165046-d7d8-4de8-ad21-5b52476e3537#testing)

## Project Structure

```
project/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── interface.go
│   │   ├── sqlite.go
│   │   ├── postgres.go
│   │   └── mysql.go
│   ├── models/
│   │   ├── product.go
│   │   └── user.go
│   ├── repository/
│   │   ├── interfaces/
│   │   │   ├── product.go
│   │   │   └── user.go
│   │   └── implementation/
│   │       ├── product.go
│   │       └── user.go
│   ├── services/
│   │   ├── product.go
│   │   └── user.go
│   └── handlers/
│       ├── product.go
│       └── user.go
├── configs/
│   ├── app.yaml
│   ├── database.yaml
│   └── constants.yaml
├── migrations/
├── go.mod
└── go.sum
```

## Dependencies

```go
// go.mod
module your-project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    gorm.io/driver/sqlite v1.5.4
    gorm.io/driver/postgres v1.5.4
    gorm.io/driver/mysql v1.5.2
    gopkg.in/yaml.v3 v3.0.1
    github.com/joho/godotenv v1.4.0
)
```

Install dependencies:

```bash
go mod init your-project
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get gorm.io/driver/postgres
go get gorm.io/driver/mysql
go get gopkg.in/yaml.v3
go get github.com/joho/godotenv
```

## Configuration

### configs/app.yaml

```yaml
server:
  host: "localhost"
  port: 8080
  mode: "debug" # debug, release, test
  timeout: 30

logging:
  level: "info"
  format: "json"
  output: "stdout"

cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["*"]
```

### configs/database.yaml

```yaml
database:
  driver: "sqlite" # sqlite, postgres, mysql

  sqlite:
    path: "data/app.db"

  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    dbname: "myapp"
    sslmode: "disable"
    timezone: "UTC"

  mysql:
    host: "localhost"
    port: 3306
    user: "root"
    password: "password"
    dbname: "myapp"
    charset: "utf8mb4"
    parse_time: true
    loc: "Local"

connection_pool:
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600 # seconds
```

### configs/constants.yaml

```yaml
pagination:
  default_page_size: 10
  max_page_size: 100

validation:
  min_password_length: 8
  max_name_length: 100

business_rules:
  max_products_per_user: 1000
  default_product_status: "active"
```

### internal/config/config.go

```go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

type Config struct {
    App       AppConfig      `yaml:"app"`
    Database  DatabaseConfig `yaml:"database"`
    Constants Constants      `yaml:"constants"`
}

type AppConfig struct {
    Server ServerConfig `yaml:"server"`
    Logging LoggingConfig `yaml:"logging"`
    CORS   CORSConfig   `yaml:"cors"`
}

type ServerConfig struct {
    Host    string `yaml:"host"`
    Port    int    `yaml:"port"`
    Mode    string `yaml:"mode"`
    Timeout int    `yaml:"timeout"`
}

type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
    Output string `yaml:"output"`
}

type CORSConfig struct {
    AllowedOrigins []string `yaml:"allowed_origins"`
    AllowedMethods []string `yaml:"allowed_methods"`
    AllowedHeaders []string `yaml:"allowed_headers"`
}

type DatabaseConfig struct {
    Driver         string               `yaml:"driver"`
    SQLite         SQLiteConfig         `yaml:"sqlite"`
    Postgres       PostgresConfig       `yaml:"postgres"`
    MySQL          MySQLConfig          `yaml:"mysql"`
    ConnectionPool ConnectionPoolConfig `yaml:"connection_pool"`
}

type SQLiteConfig struct {
    Path string `yaml:"path"`
}

type PostgresConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
    SSLMode  string `yaml:"sslmode"`
    TimeZone string `yaml:"timezone"`
}

type MySQLConfig struct {
    Host      string `yaml:"host"`
    Port      int    `yaml:"port"`
    User      string `yaml:"user"`
    Password  string `yaml:"password"`
    DBName    string `yaml:"dbname"`
    Charset   string `yaml:"charset"`
    ParseTime bool   `yaml:"parse_time"`
    Loc       string `yaml:"loc"`
}

type ConnectionPoolConfig struct {
    MaxIdleConns    int `yaml:"max_idle_conns"`
    MaxOpenConns    int `yaml:"max_open_conns"`
    ConnMaxLifetime int `yaml:"conn_max_lifetime"`
}

type Constants struct {
    Pagination    PaginationConfig    `yaml:"pagination"`
    Validation    ValidationConfig    `yaml:"validation"`
    BusinessRules BusinessRulesConfig `yaml:"business_rules"`
}

type PaginationConfig struct {
    DefaultPageSize int `yaml:"default_page_size"`
    MaxPageSize     int `yaml:"max_page_size"`
}

type ValidationConfig struct {
    MinPasswordLength int `yaml:"min_password_length"`
    MaxNameLength     int `yaml:"max_name_length"`
}

type BusinessRulesConfig struct {
    MaxProductsPerUser    int    `yaml:"max_products_per_user"`
    DefaultProductStatus  string `yaml:"default_product_status"`
}

func Load() (*Config, error) {
    config := &Config{}

    // Load app config
    if err := loadYAMLFile("configs/app.yaml", &config.App); err != nil {
        return nil, fmt.Errorf("failed to load app config: %w", err)
    }

    // Load database config
    if err := loadYAMLFile("configs/database.yaml", &config.Database); err != nil {
        return nil, fmt.Errorf("failed to load database config: %w", err)
    }

    // Load constants
    if err := loadYAMLFile("configs/constants.yaml", &config.Constants); err != nil {
        return nil, fmt.Errorf("failed to load constants: %w", err)
    }

    return config, nil
}

func loadYAMLFile(filename string, out interface{}) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    return yaml.Unmarshal(data, out)
}

func (c *Config) GetAddress() string {
    return fmt.Sprintf("%s:%d", c.App.Server.Host, c.App.Server.Port)
}
```

## Database Interface

### internal/database/interface.go

```go
package database

import (
    "gorm.io/gorm"
    "your-project/internal/config"
)

type Database interface {
    Connect() (*gorm.DB, error)
    GetDSN() string
    Close() error
}

type Manager struct {
    db     *gorm.DB
    config *config.DatabaseConfig
}

func NewManager(cfg *config.DatabaseConfig) (Database, error) {
    var db Database

    switch cfg.Driver {
    case "sqlite":
        db = NewSQLiteDB(cfg)
    case "postgres":
        db = NewPostgresDB(cfg)
    case "mysql":
        db = NewMySQLDB(cfg)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
    }

    return db, nil
}
```

### internal/database/sqlite.go

```go
package database

import (
    "fmt"
    "time"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "your-project/internal/config"
)

type SQLiteDB struct {
    config *config.DatabaseConfig
    db     *gorm.DB
}

func NewSQLiteDB(cfg *config.DatabaseConfig) *SQLiteDB {
    return &SQLiteDB{config: cfg}
}

func (s *SQLiteDB) Connect() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(s.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }

    sqlDB.SetMaxIdleConns(s.config.ConnectionPool.MaxIdleConns)
    sqlDB.SetMaxOpenConns(s.config.ConnectionPool.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(s.config.ConnectionPool.ConnMaxLifetime) * time.Second)

    s.db = db
    return db, nil
}

func (s *SQLiteDB) GetDSN() string {
    return s.config.SQLite.Path
}

func (s *SQLiteDB) Close() error {
    if s.db != nil {
        sqlDB, err := s.db.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}
```

### internal/database/postgres.go

```go
package database

import (
    "fmt"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "your-project/internal/config"
)

type PostgresDB struct {
    config *config.DatabaseConfig
    db     *gorm.DB
}

func NewPostgresDB(cfg *config.DatabaseConfig) *PostgresDB {
    return &PostgresDB{config: cfg}
}

func (p *PostgresDB) Connect() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(p.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }

    sqlDB.SetMaxIdleConns(p.config.ConnectionPool.MaxIdleConns)
    sqlDB.SetMaxOpenConns(p.config.ConnectionPool.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(p.config.ConnectionPool.ConnMaxLifetime) * time.Second)

    p.db = db
    return db, nil
}

func (p *PostgresDB) GetDSN() string {
    cfg := p.config.Postgres
    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
        cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone)
}

func (p *PostgresDB) Close() error {
    if p.db != nil {
        sqlDB, err := p.db.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}
```

## Models

### internal/models/product.go

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null"`
    Name        string         `json:"name" gorm:"size:100;not null"`
    Description string         `json:"description" gorm:"type:text"`
    Price       uint           `json:"price" gorm:"not null"`
    Status      string         `json:"status" gorm:"size:20;default:'active'"`
    UserID      uint           `json:"user_id" gorm:"not null"`
    User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type ProductFilter struct {
    Code   string `json:"code,omitempty"`
    Name   string `json:"name,omitempty"`
    Status string `json:"status,omitempty"`
    UserID uint   `json:"user_id,omitempty"`
}

type PaginationParams struct {
    Page     int `json:"page" form:"page"`
    PageSize int `json:"page_size" form:"page_size"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
}
```

### internal/models/user.go

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    Name      string         `json:"name" gorm:"size:100;not null"`
    Products  []Product      `json:"products,omitempty" gorm:"foreignKey:UserID"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

## Repository Layer

### internal/repository/interfaces/user.go

```go
package interfaces

import (
    "context"
    "your-project/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    List(ctx context.Context, filter *models.UserFilter, pagination *models.PaginationParams) ([]*models.User, int64, error)
    Update(ctx context.Context, id uint, updates map[string]interface{}) error
    Delete(ctx context.Context, id uint) error
}
```

### internal/repository/interfaces/product.go

```go
package interfaces

import (
    "context"
    "your-project/internal/models"
)

type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) error
    GetByID(ctx context.Context, id uint) (*models.Product, error)
    GetByCode(ctx context.Context, code string) (*models.Product, error)
    List(ctx context.Context, filter *models.ProductFilter, pagination *models.PaginationParams) ([]*models.Product, int64, error)
    Update(ctx context.Context, id uint, updates map[string]interface{}) error
    Delete(ctx context.Context, id uint) error
    GetByUserID(ctx context.Context, userID uint, pagination *models.PaginationParams) ([]*models.Product, int64, error)
}
```

### internal/repository/implementation/user.go

```go
package implementation

import (
    "context"
    "fmt"

    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/repository/interfaces"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) List(ctx context.Context, filter *models.UserFilter, pagination *models.PaginationParams) ([]*models.User, int64, error) {
    query := r.db.WithContext(ctx).Model(&models.User{})

    if filter != nil {
        if filter.Email != "" {
            query = query.Where("email LIKE ?", "%"+filter.Email+"%")
        }
        if filter.Name != "" {
            query = query.Where("name LIKE ?", "%"+filter.Name+"%")
        }
    }

    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if pagination != nil {
        offset := (pagination.Page - 1) * pagination.PageSize
        query = query.Offset(offset).Limit(pagination.PageSize)
    }

    var users []*models.User
    err := query.Find(&users).Error
    return users, total, err
}

func (r *userRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
    return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
```

### internal/repository/implementation/product.go

```go
package implementation

import (
    "context"
    "fmt"

    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/repository/interfaces"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).Preload("User").First(&product, id).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *productRepository) GetByCode(ctx context.Context, code string) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).Preload("User").Where("code = ?", code).First(&product).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}

func (r *productRepository) List(ctx context.Context, filter *models.ProductFilter, pagination *models.PaginationParams) ([]*models.Product, int64, error) {
    query := r.db.WithContext(ctx).Model(&models.Product{}).Preload("User")

    // Apply filters
    if filter != nil {
        if filter.Code != "" {
            query = query.Where("code LIKE ?", "%"+filter.Code+"%")
        }
        if filter.Name != "" {
            query = query.Where("name LIKE ?", "%"+filter.Name+"%")
        }
        if filter.Status != "" {
            query = query.Where("status = ?", filter.Status)
        }
        if filter.UserID != 0 {
            query = query.Where("user_id = ?", filter.UserID)
        }
    }

    // Count total records
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    if pagination != nil {
        offset := (pagination.Page - 1) * pagination.PageSize
        query = query.Offset(offset).Limit(pagination.PageSize)
    }

    var products []*models.Product
    err := query.Find(&products).Error
    return products, total, err
}

func (r *productRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
    return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *productRepository) GetByUserID(ctx context.Context, userID uint, pagination *models.PaginationParams) ([]*models.Product, int64, error) {
    query := r.db.WithContext(ctx).Model(&models.Product{}).Where("user_id = ?", userID)

    // Count total records
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    if pagination != nil {
        offset := (pagination.Page - 1) * pagination.PageSize
        query = query.Offset(offset).Limit(pagination.PageSize)
    }

    var products []*models.Product
    err := query.Find(&products).Error
    return products, total, err
}
```

## Service Layer

### internal/services/user.go

```go
package services

import (
    "context"
    "errors"
    "fmt"

    "gorm.io/gorm"
    "your-project/internal/config"
    "your-project/internal/models"
    "your-project/internal/repository/interfaces"
)

type UserService struct {
    repo      interfaces.UserRepository
    constants *config.Constants
}

func NewUserService(repo interfaces.UserRepository, constants *config.Constants) *UserService {
    return &UserService{
        repo:      repo,
        constants: constants,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
    existingUser, err := s.repo.GetByEmail(ctx, user.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("failed to check existing user: %w", err)
    }
    if existingUser != nil {
        return fmt.Errorf("user with email %s already exists", user.Email)
    }

    return s.repo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) ListUsers(ctx context.Context, filter *models.UserFilter, pagination *models.PaginationParams) (*models.PaginatedResponse, error) {
    if pagination == nil {
        pagination = &models.PaginationParams{
            Page:     1,
            PageSize: s.constants.Pagination.DefaultPageSize,
        }
    }

    if pagination.Page < 1 {
        pagination.Page = 1
    }
    if pagination.PageSize < 1 {
        pagination.PageSize = s.constants.Pagination.DefaultPageSize
    }
    if pagination.PageSize > s.constants.Pagination.MaxPageSize {
        pagination.PageSize = s.constants.Pagination.MaxPageSize
    }

    users, total, err := s.repo.List(ctx, filter, pagination)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    totalPages := (int(total) + pagination.PageSize - 1) / pagination.PageSize

    return &models.PaginatedResponse{
        Data:       users,
        Page:       pagination.Page,
        PageSize:   pagination.PageSize,
        Total:      total,
        TotalPages: totalPages,
    }, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
    if email, ok := updates["email"]; ok {
        existingUser, err := s.repo.GetByEmail(ctx, email.(string))
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("failed to check existing user: %w", err)
        }
        if existingUser != nil && existingUser.ID != id {
            return fmt.Errorf("user with email %s already exists", email)
        }
    }

    return s.repo.Update(ctx, id, updates)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}
```

### internal/services/product.go

```go
package services

import (
    "context"
    "errors"
    "fmt"
    "math"

    "gorm.io/gorm"
    "your-project/internal/config"
    "your-project/internal/models"
    "your-project/internal/repository/interfaces"
)

type ProductService struct {
    repo      interfaces.ProductRepository
    userRepo  interfaces.UserRepository
    constants *config.Constants
}

func NewProductService(repo interfaces.ProductRepository, userRepo interfaces.UserRepository, constants *config.Constants) *ProductService {
    return &ProductService{
        repo:      repo,
        userRepo:  userRepo,
        constants: constants,
    }
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
    // Validate user exists
    _, err := s.userRepo.GetByID(ctx, product.UserID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("user not found")
        }
        return fmt.Errorf("failed to validate user: %w", err)
    }

    // Check if user has reached product limit
    userProducts, _, err := s.repo.GetByUserID(ctx, product.UserID, nil)
    if err != nil {
        return fmt.Errorf("failed to check user products: %w", err)
    }

    if len(userProducts) >= s.constants.BusinessRules.MaxProductsPerUser {
        return fmt.Errorf("user has reached maximum products limit")
    }

    // Check if product with same code exists
    existingProduct, err := s.repo.GetByCode(ctx, product.Code)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("failed to check existing product: %w", err)
    }
    if existingProduct != nil {
        return fmt.Errorf("product with code %s already exists", product.Code)
    }

    // Set default status if not provided
    if product.Status == "" {
        product.Status = s.constants.BusinessRules.DefaultProductStatus
    }

    return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *ProductService) GetProductByCode(ctx context.Context, code string) (*models.Product, error) {
    return s.repo.GetByCode(ctx, code)
}

func (s *ProductService) ListProducts(ctx context.Context, filter *models.ProductFilter, pagination *models.PaginationParams) (*models.PaginatedResponse, error) {
    // Set default pagination if not provided
    if pagination == nil {
        pagination = &models.PaginationParams{
            Page:     1,
            PageSize: s.constants.Pagination.DefaultPageSize,
        }
    }

    // Validate pagination parameters
    if pagination.Page < 1 {
        pagination.Page = 1
    }
    if pagination.PageSize < 1 {
        pagination.PageSize = s.constants.Pagination.DefaultPageSize
    }
    if pagination.PageSize > s.constants.Pagination.MaxPageSize {
        pagination.PageSize = s.constants.Pagination.MaxPageSize
    }

    products, total, err := s.repo.List(ctx, filter, pagination)
    if err != nil {
        return nil, fmt.Errorf("failed to list products: %w", err)
    }

    totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

    return &models.PaginatedResponse{
        Data:       products,
        Page:       pagination.Page,
        PageSize:   pagination.PageSize,
        Total:      total,
        TotalPages: totalPages,
    }, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, updates map[string]interface{}) error {
    // Check if product exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("product not found")
        }
        return fmt.Errorf("failed to get product: %w", err)
    }

    // If updating code, check for duplicates
    if code, ok := updates["code"]; ok {
        existingProduct, err := s.repo.GetByCode(ctx, code.(string))
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("failed to check existing product: %w", err)
        }
        if existingProduct != nil && existingProduct.ID != id {
            return fmt.Errorf("product with code %s already exists", code)
        }
    }

    return s.repo.Update(ctx, id, updates)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
    // Check if product exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("product not found")
        }
        return fmt.Errorf("failed to get product: %w", err)
    }

    return s.repo.Delete(ctx, id)
}
```

## API Handlers

### internal/handlers/user.go

```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type UserHandler struct {
    service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.CreateUser(c.Request.Context(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    user, err := h.service.GetUser(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    var filter models.UserFilter
    var pagination models.PaginationParams

    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.ShouldBindQuery(&pagination); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.service.ListUsers(c.Request.Context(), &filter, &pagination)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.UpdateUser(c.Request.Context(), uint(id), updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        users.POST("", h.CreateUser)
        users.GET("", h.ListUsers)
        users.GET("/:id", h.GetUser)
        users.PUT("/:id", h.UpdateUser)
        users.DELETE("/:id", h.DeleteUser)
    }
}
```

### internal/handlers/product.go

```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type ProductHandler struct {
    service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.CreateProduct(c.Request.Context(), &product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    product, err := h.service.GetProduct(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
    var filter models.ProductFilter
    var pagination models.PaginationParams

    // Bind query parameters
    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := c.ShouldBindQuery(&pagination); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.service.ListProducts(c.Request.Context(), &filter, &pagination)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.UpdateProduct(c.Request.Context(), uint(id), updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    if err := h.service.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// RegisterRoutes registers all product routes
func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
    products := router.Group("/products")
    {
        products.POST("", h.CreateProduct)
        products.GET("", h.ListProducts)
        products.GET("/:id", h.GetProduct)
        products.PUT("/:id", h.UpdateProduct)
        products.DELETE("/:id", h.DeleteProduct)
    }
}
```

## Main Application

### cmd/server/main.go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "your-project/internal/config"
    "your-project/internal/database"
    "your-project/internal/handlers"
    "your-project/internal/models"
    "your-project/internal/repository/implementation"
    "your-project/internal/services"
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
```

## Usage Examples

### Starting the Application

```bash
# Development
go run cmd/server/main.go

# Build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

### API Examples

#### Create a User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

#### List Users

```shell
curl "http://localhost:8080/api/v1/users?page=1&page_size=5&name=Jane"
```

#### Update User

```shell
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe"
  }'
```

#### Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "code": "P001",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1200,
    "user_id": 1
  }'
```

#### List Products with Filtering and Pagination

```bash
curl "http://localhost:8080/api/v1/products?page=1&page_size=10&status=active&name=laptop"
```

#### Update a Product

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 1100,
    "description": "Updated description"
  }'
```

#### Get Product by ID

```bash
curl http://localhost:8080/api/v1/products/1
```

#### Delete a Product

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

## Database Migration

### Creating Migration Files

Create migration files in the `migrations/` directory:

```sql
-- migrations/001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- migrations/001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

```sql
-- migrations/002_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- migrations/002_create_products_table.down.sql
DROP TABLE IF EXISTS products;
```

## Testing

### Unit Tests Example

Create `internal/services/product_test.go`:

```go
package services

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "your-project/internal/config"
    "your-project/internal/models"
)

type MockProductRepository struct {
    mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByCode(ctx context.Context, code string) (*models.Product, error) {
    args := m.Called(ctx, code)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context, filter *models.ProductFilter, pagination *models.PaginationParams) ([]*models.Product, int64, error) {
    args := m.Called(ctx, filter, pagination)
    return args.Get(0).([]*models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
    args := m.Called(ctx, id, updates)
    return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockProductRepository) GetByUserID(ctx context.Context, userID uint, pagination *models.PaginationParams) ([]*models.Product, int64, error) {
    args := m.Called(ctx, userID, pagination)
    return args.Get(0).([]*models.Product), args.Get(1).(int64), args.Error(2)
}

func TestProductService_CreateProduct(t *testing.T) {
    mockRepo := new(MockProductRepository)
    mockUserRepo := new(MockUserRepository)

    constants := &config.Constants{
        BusinessRules: config.BusinessRulesConfig{
            MaxProductsPerUser:   1000,
            DefaultProductStatus: "active",
        },
    }

    service := NewProductService(mockRepo, mockUserRepo, constants)

    product := &models.Product{
        Code:   "TEST001",
        Name:   "Test Product",
        Price:  100,
        UserID: 1,
    }

    // Mock user exists
    mockUserRepo.On("GetByID", mock.Anything, uint(1)).Return(&models.User{ID: 1}, nil)

    // Mock no existing products for user
    mockRepo.On("GetByUserID", mock.Anything, uint(1), (*models.PaginationParams)(nil)).Return([]*models.Product{}, int64(0), nil)

    // Mock no existing product with same code
    mockRepo.On("GetByCode", mock.Anything, "TEST001").Return(nil, gorm.ErrRecordNotFound)

    // Mock successful creation
    mockRepo.On("Create", mock.Anything, product).Return(nil)

    err := service.CreateProduct(context.Background(), product)

    assert.NoError(t, err)
    assert.Equal(t, "active", product.Status)
    mockRepo.AssertExpectations(t)
    mockUserRepo.AssertExpectations(t)
}
```

### Integration Tests

Create `cmd/server/main_test.go`:

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "your-project/internal/config"
    "your-project/internal/database"
    "your-project/internal/handlers"
    "your-project/internal/models"
    "your-project/internal/repository/implementation"
    "your-project/internal/services"
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
```

## Key Features

### 1. Modular Architecture

- **Separation of Concerns**: Clear separation between database, repository, service, and handler layers
- **Interface-based Design**: Easy to mock and test
- **Pluggable Database Support**: Switch between SQLite, PostgreSQL, and MySQL

### 2. Configuration Management

- **YAML Configuration**: Easy to read and modify
- **Environment-specific configs**: Support for different environments
- **Constants Management**: Centralized business rules and constants

### 3. Repository Pattern Benefits

- **Database Abstraction**: Business logic doesn't depend on specific database implementation
- **Testability**: Easy to mock repositories for unit testing
- **Consistency**: Uniform data access patterns

### 4. Production-Ready Features

- **Connection Pooling**: Optimized database connections
- **Graceful Shutdown**: Proper cleanup on application termination
- **CORS Support**: Cross-origin resource sharing
- **Pagination**: Built-in pagination support
- **Error Handling**: Comprehensive error handling and validation

### 5. Extensibility

- **Easy Database Addition**: Add new database drivers by implementing the Database interface
- **Plugin Architecture**: Easy to add new features without modifying existing code
- **Service Layer**: Business logic separation allows for easy feature extension

## Switching Databases

To switch from SQLite to PostgreSQL, simply update `configs/database.yaml`:

```yaml
database:
  driver: "postgres"  # Change from "sqlite"
  # ... postgres configuration
```

The application will automatically use the appropriate database driver without any code changes.

## Best Practices

1. **Always use transactions for complex operations**
2. **Implement proper error handling and logging**
3. **Use context for request tracing and cancellation**
4. **Validate input data at the handler level**
5. **Keep business logic in the service layer**
6. **Write comprehensive tests for all layers**
7. **Use database migrations for schema changes**
8. **Monitor database connection pool metrics**
9. **Implement proper security measures (authentication, authorization)**
10. **Use environment variables for sensitive configuration**

This architecture provides a solid foundation for building scalable, maintainable Go applications with clean separation of concerns and easy database portability.

---

## Notes

Documentation for implementing a modular repository pattern with GORM, Gin, and YAML configuration. Here are the key highlights:

## What's Included:

### 🏗️ **Architecture**

- **Modular Design**: Clean separation between database, repository, service, and handler layers
- **Interface-based**: Easy to mock and test
- **Multi-database Support**: SQLite, PostgreSQL, and MySQL with easy switching

### 📁 **Project Structure**

- Organized folder structure following Go best practices
- Clear separation of concerns across different packages
- Scalable architecture for growing applications

### ⚙️ **Configuration Management**

- **YAML-based config**: Separate files for app, database, and constants
- **Environment flexibility**: Easy to switch between development, staging, and production
- **Business rules centralization**: Constants and validation rules in one place

### 🔌 **Database Abstraction**

- **Database Interface**: Common interface for all database types
- **Connection Pooling**: Optimized database connections
- **Auto-migration**: Automatic schema migration support

### 📊 **Repository Pattern**

- **Interface-driven**: Repository interfaces for easy testing
- **CRUD Operations**: Complete Create, Read, Update, Delete operations
- **Advanced Features**: Filtering, pagination, and complex queries

### 🚀 **Production-Ready Features**

- **Graceful Shutdown**: Proper cleanup on application termination
- **Health Checks**: Built-in health check endpoints
- **CORS Support**: Cross-origin resource sharing
- **Error Handling**: Comprehensive error handling and validation
- **Logging**: Structured logging support

### 🧪 **Testing Support**

- **Unit Tests**: Examples with mocking
- **Integration Tests**: Complete end-to-end testing examples
- **Test Database**: In-memory SQLite for testing

## Key Benefits:

1. **Easy Database Switching**: Change one line in YAML config to switch databases
2. **Testable**: Interface-based design makes unit testing straightforward
3. **Scalable**: Clean architecture supports growing applications
4. **Maintainable**: Clear separation of concerns and well-organized code
5. **Production-Ready**: Includes all necessary features for production deployment

## Quick Start:

1. Copy the project structure
2. Update `configs/database.yaml` with your database preferences
3. Run `go mod tidy` to install dependencies
4. Run `go run cmd/server/main.go` to start the server

The documentation includes complete working examples, API usage, testing strategies, and best practices for building robust Go applications with this architecture.

---

## Authentication & Authorization

Here's a comprehensive JWT authentication and authorization implementation with middleware to restrict user deletion to admins or account owners:

### 1. Update User Model (`internal/models/user.go`)

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    Name      string         `json:"name" gorm:"size:100;not null"`
    Password  string         `json:"-" gorm:"size:255;not null"` // Hashed password
    Role      Role           `json:"role" gorm:"type:enum('user','admin');default:'user'"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Add to constants.yaml
type AuthConfig struct {
    JWTSecret           string `yaml:"jwt_secret"`
    TokenExpiration     int    `yaml:"token_expiration"` // in hours
    PasswordCost        int    `yaml:"password_cost"`
}
```

### 2. Add Auth Configuration (`configs/constants.yaml`)

```yaml
auth:
  jwt_secret: "your-strong-secret-key" # Should be 32+ chars
  token_expiration: 72   # hours
  password_cost: 14      # bcrypt cost factor
```

### 3. JWT Service (`internal/services/auth.go`)

```go
package services

import (
    "context"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/config"
    "your-project/internal/models"
)

type AuthService struct {
    jwtSecret   string
    tokenExpiry time.Duration
    passwordCost int
}

func NewAuthService(constants *config.Constants) *AuthService {
    return &AuthService{
        jwtSecret:   constants.Auth.JWTSecret,
        tokenExpiry: time.Duration(constants.Auth.TokenExpiration) * time.Hour,
        passwordCost: constants.Auth.PasswordCost,
    }
}

type Claims struct {
    UserID uint `json:"user_id"`
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
```

### 4. Auth Middleware (`internal/middleware/auth.go`)

```go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "your-project/internal/services"
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

        // Set user context
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

        // Check if user has any of the allowed roles
        hasPermission := false
        for _, allowedRole := range allowedRoles {
            if role == allowedRole {
                hasPermission = true
                break
            }
        }

        if !hasPermission {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            return
        }

        c.Next()
    }
}
```

### 5. Update User Service (`internal/services/user.go`)

```go
// Add to UserService struct
authService *AuthService

// Update NewUserService
func NewUserService(repo interfaces.UserRepository, authService *AuthService, constants *config.Constants) *UserService {
    return &UserService{
        repo:        repo,
        authService: authService,
        constants:   constants,
    }
}

// Update CreateUser method to hash password
func (s *UserService) CreateUser(ctx context.Context, user *models.User, password string) error {
    // ... existing validation ...

    hashedPassword, err := s.authService.HashPassword(password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    user.Password = hashedPassword
    return s.repo.Create(ctx, user)
}

// Add Login method
func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
    user, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        return nil, "", err
    }

    if !s.authService.CheckPasswordHash(password, user.Password) {
        return nil, "", errors.New("invalid credentials")
    }

    token, err := s.authService.GenerateToken(user)
    if err != nil {
        return nil, "", fmt.Errorf("failed to generate token: %w", err)
    }

    return user, token, nil
}
```

### 6. Auth Handler (`internal/handlers/auth.go`)

```go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type AuthHandler struct {
    userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
    return &AuthHandler{userService: userService}
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, token, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    response := gin.H{
        "user":  user,
        "token": token,
    }

    c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
    auth := router.Group("/auth")
    {
        auth.POST("/login", h.Login)
    }
}
```

### 7. Update User Handler (`internal/handlers/user.go`)

```go
// Add to UserHandler struct
authService *services.AuthService

// Update NewUserHandler
func NewUserHandler(service *services.UserService, authService *services.AuthService) *UserHandler {
    return &UserHandler{
        service:     service,
        authService: authService,
    }
}

// Update CreateUser handler
func (h *UserHandler) CreateUser(c *gin.Context) {
    type CreateUserRequest struct {
        models.User
        Password string `json:"password" binding:"required,min=8"`
    }

    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.CreateUser(c.Request.Context(), &req.User, req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Clear password before returning
    req.User.Password = ""
    c.JSON(http.StatusCreated, req.User)
}

// Add authorization to DeleteUser
func (h *UserHandler) DeleteUser(c *gin.Context) {
    // Get user ID from token
    authUserID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
        return
    }

    // Get user role from token
    authUserRole, exists := c.Get("userRole")
    if !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "Authentication required"})
        return
    }

    // Get user ID from path
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    // Check authorization: only admins or account owners can delete
    isAdmin := authUserRole == models.RoleAdmin
    isOwner := authUserID.(uint) == uint(id)

    if !isAdmin && !isOwner {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this user"})
        return
    }

    if err := h.service.DeleteUser(c.Request.Context(), uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
```

### 8. Update Main Application (`cmd/server/main.go`)

```go
// Initialize services
authService := services.NewAuthService(&cfg.Constants)
userService := services.NewUserService(userRepo, authService, &cfg.Constants)
productService := services.NewProductService(productRepo, userRepo, &cfg.Constants)

// Initialize handlers
authHandler := handlers.NewAuthHandler(userService)
userHandler := handlers.NewUserHandler(userService, authService)
productHandler := handlers.NewProductHandler(productService)

// Setup router
router := gin.New()

// ... existing middleware ...

// Public routes
public := router.Group("/api/v1")
authHandler.RegisterRoutes(public)

// Protected routes
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(authService))

// Admin-only routes
admin := protected.Group("")
admin.Use(middleware.RoleMiddleware(models.RoleAdmin))

// Apply auth middleware to user routes
userHandler.RegisterRoutes(protected)

// Example admin-only route
admin.GET("/admin-only", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Admin access granted"})
})

// ... start server ...
```

### 9. Authorization Middleware Usage

```go
// User deletion route (already protected by AuthMiddleware)
users.DELETE("/:id", h.DeleteUser)

// Admin-only route example
admin := router.Group("/admin")
admin.Use(middleware.RoleMiddleware(models.RoleAdmin))
{
    admin.GET("/dashboard", adminDashboardHandler)
}
```

### 10. JWT Usage Example

```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "securepassword"}'

# Response: 
# {"user": {...}, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}

# Delete user (as admin or owner)
curl -X DELETE http://localhost:8080/api/v1/users/123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Admin-only endpoint
curl http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Key Security Features:

1. **JWT Authentication**:
  
  - Bearer token authentication
  - HMAC-SHA256 signing
  - Configurable expiration
  - Secure secret management
2. **Password Security**:
  
  - Bcrypt hashing with configurable cost
  - Password never stored or returned in plain text
  - Minimum password length enforcement
3. **Authorization**:
  
  - Role-based access control (RBAC)
  - Ownership verification
  - Middleware for route protection
4. **Secure Practices**:
  
  - Authorization header validation
  - Proper error messages (avoid information leakage)
  - Context propagation of user identity
  - Separation of public/protected routes

### Implementation Notes:

1. **Secrets Management**:
  
  - Store JWT secret in environment variables in production
  - Use `godotenv` to load from `.env` file in development
  - Rotate secrets regularly
2. **Password Handling**:
  
  - Always hash passwords before storage
  - Never return password hashes in API responses
  - Implement password complexity requirements
3. **Token Security**:
  
  - Use HTTPS in production
  - Set appropriate token expiration times
  - Implement token revocation for sensitive operations
4. **Authorization Logic**:
  
  - The middleware handles role-based access
  - Resource ownership is checked at the handler level
  - Combine both for fine-grained control

This implementation provides a robust authentication and authorization system that prevents users from deleting other users' accounts unless they are admins. The middleware architecture ensures these security checks are consistently applied across all protected routes.

---
