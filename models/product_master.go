package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductMaster struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductDocumentID string    `gorm:"type:char(36)" json:"product_document_id"`
	Barcode           string    `gorm:"size:255;not null" json:"barcode"`
	BarcodeWarehouse  string    `gorm:"size:255" json:"barcode_warehouse"`
	BarcodeSource     string    `gorm:"size:50" json:"barcode_source"`
	Qty               int       `json:"qty"`
	QtySource         int       `json:"qty_source"`
	Price             float64   `gorm:"type:decimal(15,2)" json:"price"`
	PriceSource       float64   `gorm:"type:decimal(15,2)" json:"price_source"`
	PriceWarehouse    float64   `gorm:"type:decimal(15,2)" json:"price_warehouse"`
	CategoryID        string    `gorm:"type:char(36)" json:"category_id"`
	StickerID         string    `gorm:"type:char(36)" json:"sticker_id"`
	Reguleher         string    `gorm:"size:50" json:"reguleher"`
	OutQuantity       int       `gorm:"default:0" json:"out_quantity"`
	IsReidentify      bool      `gorm:"default:false" json:"is_reidentify"`
	BundleID          string    `gorm:"type:char(36)" json:"bundle_id"`
	SKUId             string    `gorm:"type:char(36)" json:"sku_id"`
	BagID             string    `gorm:"type:char(36)" json:"bag_id"`
	UserID            string    `gorm:"type:char(36)" json:"user_id"`
	Status            string    `gorm:"size:50" json:"status"`
	StatusSource      string    `gorm:"size:50" json:"status_source"`
	StoreID           int       `json:"store_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
