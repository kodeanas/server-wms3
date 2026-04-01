package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreCrew struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:255" json:"email"`
	Address   string    `gorm:"type:text" json:"address"`
	StoreID   string    `gorm:"type:char(36);not null" json:"store_id"`
	IsCashier bool      `gorm:"default:false" json:"is_cashier"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
