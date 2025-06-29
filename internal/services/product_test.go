package services

import (
	"context"
	"testing"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter *models.User, pagination *models.PaginationParams) ([]*models.User, int64, error) {
    args := m.Called(ctx, filter, pagination)
    return args.Get(0).([]*models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
    args := m.Called(ctx, id, updates)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

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
