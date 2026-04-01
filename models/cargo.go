package models

import (
	"time"

	"github.com/google/uuid"
)

type Cargo struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code          string    `gorm:"size:255;not null" json:"code"`
	Datetime      time.Time `gorm:"autoCreateTime" json:"datetime"`
	UserID        string    `gorm:"type:char(36)" json:"user_id"`
	TotalQuantity int       `json:"total_quantity"`
	TotalPrice    float64   `gorm:"type:decimal(15,2)" json:"total_price"`
	Status        string    `gorm:"size:50" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
