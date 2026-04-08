package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductPending struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DocumentID  string         `gorm:"type:char(36)" json:"document_id"`
	Barcode     string         `gorm:"size:255" json:"barcode"`
	Name        string         `gorm:"type:text" json:"name"`
	Item        int            `json:"item"`
	Price       float64        `gorm:"type:decimal(15,2)" json:"price"`
	IsSKU       bool           `gorm:"default:false" json:"is_sku"`
	Status      string         `gorm:"size:50" json:"status"`
	Note        string         `gorm:"type:text" json:"note"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DateScanned *time.Time     `json:"date_scanned"`
}
