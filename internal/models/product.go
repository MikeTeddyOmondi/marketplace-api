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
