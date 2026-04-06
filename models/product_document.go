package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductDocument struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code          string         `gorm:"size:255;not null" json:"code"`
	FileName      string         `gorm:"size:255" json:"file_name"`
	FileItem      int            `json:"file_item"`
	FilePrice     int            `json:"file_price"`
	Status        string         `gorm:"size:50" json:"status"`
	Type          string         `gorm:"size:50" json:"type"`
	HeaderBarcode string         `gorm:"size:255" json:"header_barcode"`
	HeaderName    string         `gorm:"size:255" json:"header_name"`
	HeaderItem    string         `gorm:"size:255" json:"header_item"`
	HeaderPrice   string         `gorm:"size:255" json:"header_price"`
	UserID        *string        `gorm:"type:char(36)" json:"user_id"`
	Supplier      string         `gorm:"size:255" json:"supplier"`
	TypeProduct   string         `gorm:"size:50" json:"type_product"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
