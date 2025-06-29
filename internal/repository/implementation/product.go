package implementation

import (
	"context"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/interfaces"

	"gorm.io/gorm"
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
