package models

import (
	"time"

	"github.com/google/uuid"
)

type StoreTransfer struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StoreID             string    `gorm:"type:char(36)" json:"store_id"`
	Datetime            time.Time `gorm:"autoCreateTime" json:"datetime"`
	TotalItem           int       `json:"total_item"`
	TotalPrice          float64   `gorm:"type:decimal(15,2)" json:"total_price"`
	TotalPriceWarehouse float64   `gorm:"type:decimal(15,2)" json:"total_price_warehouse"`
	Status              string    `gorm:"size:50" json:"status"`
	UserID              string    `gorm:"type:char(36)" json:"user_id"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
