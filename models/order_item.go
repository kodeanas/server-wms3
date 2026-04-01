package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductMasterID string    `gorm:"type:char(36)" json:"product_master_id"`
	OrderID         string    `gorm:"type:char(36)" json:"order_id"`
	Price           float64   `gorm:"type:decimal(15,2)" json:"price"`
	PriceWarehouse  float64   `gorm:"type:decimal(15,2)" json:"price_warehouse"`
	PriceCut        float64   `gorm:"type:decimal(15,2)" json:"price_cut"`
	PrceFinal       float64   `gorm:"type:decimal(15,2)" json:"prce_final"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
