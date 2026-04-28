package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductOrder struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID        string         `gorm:"type:char(36)" json:"order_id"`
	ProductID      string         `gorm:"type:char(36)" json:"product_id"`
	Name           string         `gorm:"size:255" json:"name"`
	Price          float64        `gorm:"type:decimal(15,2)" json:"price"`
	PriceWarehouse float64        `gorm:"type:decimal(15,2)" json:"price_warehouse"`
	Discount       float64        `gorm:"type:decimal(15,2)" json:"discount"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
