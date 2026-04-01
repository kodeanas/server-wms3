package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductDocument struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Code           string    `gorm:"size:255;not null" json:"code"`
	HeaderBarcode  string    `gorm:"size:255" json:"header_barcode"`
	HeaderName     string    `gorm:"size:255" json:"header_name"`
	HeaderQuantity string    `gorm:"size:255" json:"header_quantity"`
	HeaderPrice    string    `gorm:"size:255" json:"header_price"`
	TypeIn         string    `gorm:"size:50" json:"type_in"`
	TypeBulking    string    `gorm:"size:50" json:"type_bulking"`
	TotalList      int       `json:"total_list"`
	TotalPrice     float64   `gorm:"type:decimal(15,2)" json:"total_price"`
	UserID         string    `gorm:"type:char(36)" json:"user_id"`
	SupplierID     string    `gorm:"type:char(36)" json:"supplier_id"`
	Status         string    `gorm:"size:50;default:'active'" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
