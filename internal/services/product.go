package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/interfaces"

	"gorm.io/gorm"
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
