package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sticker struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CodeHex    string         `gorm:"size:255;not null" json:"code_hex"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Slug       string         `gorm:"size:255;unique;not null" json:"slug"`
	Type       string         `gorm:"size:50" json:"type"`
	FixedPrice *int           `json:"fixed_price"`
	MinPrice   *Price         `gorm:"type:decimal(15,2)" json:"min_price"`
	MaxPrice   *Price         `gorm:"type:decimal(15,2)" json:"max_price"`
	Status     string         `gorm:"size:50;default:'active'" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// CreateStickerPayload request payload.
type CreateStickerPayload struct {
	CodeHex    string   `json:"code_hex" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Status     string   `json:"status"`
}

// UpdateStickerPayload request payload for update.
type UpdateStickerPayload struct {
	CodeHex    string   `json:"code_hex"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Status     string   `json:"status"`
}
