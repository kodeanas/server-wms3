package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SlowMoving represents slow moving product tracking
type SlowMoving struct {
	ID                  string    `gorm:"primaryKey;type:uuid" json:"id"`
	Date                time.Time `json:"date"`
	TotalItem           int       `json:"total_item"`
	TotalPrice          string    `gorm:"type:numeric(19,2)" json:"total_price"`
	TotalPriceWarehouse string    `gorm:"type:numeric(19,2)" json:"total_price_warehouse"`
	IsDamaged           bool      `json:"is_damaged"`
	UserID              string    `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relationships
	User  *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []SlowMovingItem `gorm:"foreignKey:SlowMovingID" json:"items,omitempty"`
}

// BeforeCreate generates UUID before creating
func (sm *SlowMoving) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == "" {
		sm.ID = uuid.New().String()
	}
	return nil
}

// SlowMovingItem represents items in slow moving tracking
type SlowMovingItem struct {
	ID              string    `gorm:"primaryKey;type:uuid" json:"id"`
	SlowMovingID    string    `gorm:"type:uuid;not null;index" json:"slow_moving_id"`
	ProductMasterID string    `gorm:"type:uuid;not null;index" json:"product_master_id"`
	IsDamaged       bool      `json:"is_damaged"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	SlowMoving    *SlowMoving    `gorm:"foreignKey:SlowMovingID" json:"slow_moving,omitempty"`
	ProductMaster *ProductMaster `gorm:"foreignKey:ProductMasterID" json:"product_master,omitempty"`
}

// BeforeCreate generates UUID before creating
func (smi *SlowMovingItem) BeforeCreate(tx *gorm.DB) error {
	if smi.ID == "" {
		smi.ID = uuid.New().String()
	}
	return nil
}
