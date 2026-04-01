package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Slug      string    `gorm:"size:255;unique;not null" json:"slug"`
	Discount  int       `gorm:"default:0" json:"discount"`
	MinPrice  float64   `gorm:"type:decimal(15,2)" json:"min_price"`
	MaxPrice  float64   `gorm:"type:decimal(15,2)" json:"max_price"`
	Status    string    `gorm:"size:50;default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
