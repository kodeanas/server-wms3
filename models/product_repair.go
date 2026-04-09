package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepair struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID    *uuid.UUID     `gorm:"type:uuid" json:"product_id"`
	DocumentID   *uuid.UUID     `gorm:"type:uuid" json:"document_id"`
	Status       string         `gorm:"size:50" json:"status"`
	DateOut      *time.Time     `json:"date_out"`
	UserID       *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	PriceBefore  float64        `gorm:"type:decimal(15,2)" json:"price_before"`
	PriceUpdate  float64        `gorm:"type:decimal(15,2)" json:"price_update"`
	ItemBefore   int            `json:"item_before"`
	ItemUpdate   int            `json:"item_update"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	StickerID    *uuid.UUID     `gorm:"type:uuid" json:"sticker_id"`
	RemarkOrigin string         `gorm:"type:text" json:"remark_origin"`
	RemarkAfter  string         `gorm:"type:text" json:"remark_after"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
