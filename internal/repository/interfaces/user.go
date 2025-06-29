package interfaces

import (
	"context"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, filter *models.User, pagination *models.PaginationParams) ([]*models.User, int64, error)
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
}
