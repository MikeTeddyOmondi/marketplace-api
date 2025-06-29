package interfaces

import (
	"context"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
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
