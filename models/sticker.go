package models

import (
	"time"

	"github.com/google/uuid"
)

type Sticker struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CodeHex    string    `gorm:"size:255;not null" json:"code_hex"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Slug       string    `gorm:"size:255;unique;not null" json:"slug"`
	Type       string    `gorm:"size:50" json:"type"`
	FixedPrice int       `json:"fixed_price"`
	MinPrice   float64   `gorm:"type:decimal(15,2)" json:"min_price"`
	MaxPrice   float64   `gorm:"type:decimal(15,2)" json:"max_price"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
