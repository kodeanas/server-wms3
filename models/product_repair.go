package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepair struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DocumentID  string         `gorm:"type:char(36)" json:"document_id"`
	Barcode     string         `gorm:"size:255" json:"barcode"`
	Name        string         `gorm:"type:text" json:"name"`
	ItemBefore  int            `json:"item_before"`
	ItemUpdate  int            `json:"item_update"`
	PriceBefore float64        `gorm:"type:decimal(15,2)" json:"price_before"`
	PriceUpdate float64        `gorm:"type:decimal(15,2)" json:"price_update"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
