package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreTransferBag struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StoreTransferID string    `gorm:"type:char(36)" json:"store_transfer_id"`
	BagID           string    `gorm:"type:char(36)" json:"bag_id"`
	Quantity        int       `json:"quantity"`
	TotalPrice      float64   `gorm:"type:decimal(15,2)" json:"total_price"`
	TotalCOGS       float64   `gorm:"type:decimal(15,2)" json:"total_cogs"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
