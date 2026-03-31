package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductDocument represents a product document
type ProductDocument struct {
	ID             string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name           string    `gorm:"type:varchar(255);not null" json:"name"`
	Code           string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"code"`
	HeaderBarcode  string    `gorm:"type:varchar(255)" json:"header_barcode"`
	HeaderName     string    `gorm:"type:varchar(255)" json:"header_name"`
	HeaderQuantity string    `gorm:"type:varchar(255)" json:"header_quantity"`
	HeaderPrice    string    `gorm:"type:varchar(255)" json:"header_price"`
	TypeIn         string    `gorm:"type:varchar(50)" json:"type_in"`      // enum
	TypeBulking    string    `gorm:"type:varchar(50)" json:"type_bulking"` // enum
	UserID         int       `json:"user_id"`
	SupplierID     string    `gorm:"type:uuid" json:"supplier_id"`
	TotalList      int       `json:"total_list"`
	TotalPrice     string    `gorm:"type:numeric(19,2)" json:"total_price"`
	Status         string    `gorm:"type:varchar(50);default:'done'" json:"status"` // enum
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	ProductMasters []ProductMaster `gorm:"foreignKey:ProductDocumentID" json:"product_masters,omitempty"`
}

// BeforeCreate generates UUID before creating
func (pd *ProductDocument) BeforeCreate(tx *gorm.DB) error {
	if pd.ID == "" {
		pd.ID = uuid.New().String()
	}
	return nil
}

// ProductMaster represents a product master record
type ProductMaster struct {
	ID                string         `gorm:"primaryKey;type:uuid" json:"id"`
	ProductDocumentID string         `gorm:"type:uuid;index" json:"product_document_id"`
	Barcode           string         `gorm:"type:varchar(255);not null" json:"barcode"`
	BarcodeWarehouse  string         `gorm:"type:varchar(255)" json:"barcode_warehouse"`
	Qty               int            `gorm:"not null" json:"qty"`
	QtySource         int            `json:"qty_source"` // enum
	Price             string         `gorm:"type:numeric(19,2);not null" json:"price"`
	PriceSource       string         `gorm:"type:numeric(19,2)" json:"price_source"` // enum
	PriceWarehouse    string         `gorm:"type:numeric(19,2)" json:"price_warehouse"`
	CategoryID        string         `gorm:"type:uuid;index" json:"category_id"`
	StickerID         string         `gorm:"type:uuid;index" json:"sticker_id"`
	BundleID          string         `gorm:"type:uuid" json:"bundle_id"`
	SKUID             string         `gorm:"type:uuid" json:"sku_id"`
	BagID             string         `gorm:"type:uuid" json:"bag_id"`
	UserID            string         `gorm:"type:uuid;index" json:"user_id"`
	Out               int            `gorm:"type:varchar(50)" json:"out"` // enum: regular, wholesale
	IsReidentify      bool           `json:"is_reidentify"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DateTime          time.Time      `json:"date_time"`
	RackID            sql.NullString `gorm:"type:uuid" json:"rack_id"`

	// Relationships
	ProductDocument *ProductDocument `gorm:"foreignKey:ProductDocumentID" json:"product_document,omitempty"`
	Category        *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Sticker         *Sticker         `gorm:"foreignKey:StickerID" json:"sticker,omitempty"`
	User            *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Rack            *Rack            `gorm:"foreignKey:RackID" json:"rack,omitempty"`
	ProductLogs     []ProductLog     `gorm:"foreignKey:ProductMasterID" json:"product_logs,omitempty"`
}

// BeforeCreate generates UUID before creating
func (pm *ProductMaster) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == "" {
		pm.ID = uuid.New().String()
	}
	return nil
}

// ProductLog represents product activity log
type ProductLog struct {
	ID              string    `gorm:"primaryKey;type:uuid" json:"id"`
	ProductMasterID string    `gorm:"type:uuid;not null;index" json:"product_master_id"`
	PrevData        string    `gorm:"type:text" json:"prev_data"`
	NewData         string    `gorm:"type:text" json:"new_data"`
	DateTime        time.Time `json:"date_time"`
	UserID          string    `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`

	// Relationships
	ProductMaster *ProductMaster `gorm:"foreignKey:ProductMasterID" json:"product_master,omitempty"`
	User          *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate generates UUID before creating
func (pl *ProductLog) BeforeCreate(tx *gorm.DB) error {
	if pl.ID == "" {
		pl.ID = uuid.New().String()
	}
	return nil
}
