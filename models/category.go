package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Slug      string         `gorm:"size:255;unique;not null" json:"slug"`
	Discount  *int           `gorm:"default:0" json:"discount"`
	MinPrice  *Price         `gorm:"type:decimal(15,2)" json:"min_price"`
	MaxPrice  *Price         `gorm:"type:decimal(15,2)" json:"max_price"`
	Status    string         `gorm:"size:50;default:'active'" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// UpdateCategoryPayload request payload for update.
type UpdateCategoryPayload struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Discount *int     `json:"discount"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Status   string   `json:"status"`
}

// CreateCategoryPayload request payload.
type CreateCategoryPayload struct {
	Name     string   `json:"name" binding:"required"`
	Slug     string   `json:"slug"`
	Discount *int     `json:"discount"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
}
