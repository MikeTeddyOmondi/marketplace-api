package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/models"
	"github.com/MikeTeddyOmondi/marketplace-api/internal/repository/interfaces"

	"gorm.io/gorm"
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

func (s *UserService) ListUsers(ctx context.Context, filter *models.User, pagination *models.PaginationParams) (*models.PaginatedResponse, error) {
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
