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
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Name     string `json:"name" gorm:"size:100;not null"`
	Password string `json:"-" gorm:"size:255;not null"`                           // Hashed password
	Role     Role   `json:"role" gorm:"type:enum('user','admin');default:'user'"` // Support for MySQL & Postgres
	// Role      Role           `json:"role" gorm:"type:text;default:'user'"` // Support for SQLite
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
