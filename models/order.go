package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	IsTax               bool           `gorm:"default:false" json:"is_tax"`
	ID                  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code                string         `gorm:"size:255;not null" json:"code"`
	UserID              string         `gorm:"type:char(36)" json:"user_id"`
	Type                string         `gorm:"size:50" json:"type"`
	BuyerID             string         `gorm:"type:char(36)" json:"buyer_id"`
	Datetime            time.Time      `gorm:"autoCreateTime" json:"datetime"`
	TotalItem           int            `json:"total_item"`
	TotalPriceProduct   float64        `gorm:"type:decimal(15,2)" json:"total_price_product"`
	TotalPPN            float64        `gorm:"type:decimal(15,2)" json:"total_ppn"`
	CartonBoxPrice      float64        `gorm:"type:decimal(15,2)" json:"carton_box_price"`
	QuantityCartonBox   int            `json:"quantity_carton_box"`
	Voucher             float64        `gorm:"type:decimal(15,2)" json:"voucher"`
	FixedDiscount       int            `json:"fixed_discount"`
	ClassDiscount       int            `json:"class_discount"`
	ClassDiscountAmount float64        `gorm:"type:decimal(15,2)" json:"class_discount_amount"`
	ClassID             string         `gorm:"type:char(36)" json:"class_id"`
	GrandTotalBefore    float64        `gorm:"type:decimal(15,2)" json:"grand_total_before"`
	GrandTotalAfter     float64        `gorm:"type:decimal(15,2)" json:"grand_total_after"`
	Status              string         `gorm:"size:50" json:"status"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
