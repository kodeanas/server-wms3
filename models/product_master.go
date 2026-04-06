package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductMaster struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DocumentID       string         `gorm:"type:char(36)" json:"document_id"`
	Barcode          string         `gorm:"size:255" json:"barcode"`
	BarcodeWarehouse string         `gorm:"size:255" json:"barcode_warehouse"`
	Name             string         `gorm:"type:text" json:"name"`
	NameWarehouse    string         `gorm:"type:text" json:"name_warehouse"`
	Item             int            `json:"item"`
	ItemWarehouse    int            `json:"item_warehouse"`
	Price            float64        `gorm:"type:decimal(15,2)" json:"price"`
	PriceWarehouse   float64        `gorm:"type:decimal(15,2)" json:"price_warehouse"`
	CategoryID       *string        `gorm:"type:char(36)" json:"category_id"`
	StickerID        *string        `gorm:"type:char(36)" json:"sticker_id"`
	ProductPendingID *string        `gorm:"type:char(36)" json:"product_pending_id"`
	TypeID           string         `json:"type_id" gorm:"-"` // hanya untuk response, tidak masuk DB
	IsSKU            bool           `gorm:"default:false" json:"is_sku"`
	Location         string         `gorm:"size:50" json:"location"`
	BundleParentID   *string        `gorm:"type:char(36)" json:"bundle_parent_id"`
	DateOut          *time.Time     `json:"date_out"`
	TypeOut          string         `gorm:"size:50" json:"type_out"`
	RackStagingID    *string        `gorm:"type:char(36)" json:"rack_staging_id"`
	RackDisplayID    *string        `gorm:"type:char(36)" json:"rack_display_id"`
	BagID            *string        `gorm:"type:char(36)" json:"bag_id"`
	UserID           *string        `gorm:"type:char(36)" json:"user_id"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
