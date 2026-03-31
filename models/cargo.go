package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Cargo represents a cargo in the warehouse
type Cargo struct {
	ID            string         `gorm:"primaryKey;type:uuid" json:"id"`
	Code          string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"code"`
	DateTime      time.Time      `json:"date_time"`
	UserID        string         `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalQuantity sql.NullInt64  `json:"total_quantity"`
	TotalPrice    sql.NullString `gorm:"type:numeric(19,2)" json:"total_price"`
	Status        string         `gorm:"type:varchar(50);default:'on_progress'" json:"status"` // enum: on_progress, done
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	// Relationships
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderCargos []OrderCargo `gorm:"foreignKey:CargoID" json:"order_cargos,omitempty"`
}

// BeforeCreate generates UUID before creating
func (c *Cargo) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Bag represents a product bag/package
type Bag struct {
	ID            string         `gorm:"primaryKey;type:uuid" json:"id"`
	CargoID       sql.NullString `gorm:"type:uuid;index" json:"cargo_id"`
	Code          string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"code"`
	Type          string         `gorm:"type:varchar(50)" json:"type"` // enum: regular, sku, sticker, hl, scrap, qcd
	DateTime      time.Time      `json:"date_time"`
	UserID        string         `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalQuantity sql.NullInt64  `json:"total_quantity"`
	TotalPrice    sql.NullString `gorm:"type:numeric(19,2)" json:"total_price"`
	Status        string         `gorm:"type:varchar(50);default:'done'" json:"status"` // enum
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	// Relationships
	Cargo             *Cargo             `gorm:"foreignKey:CargoID" json:"cargo,omitempty"`
	User              *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Products          []ProductMaster    `gorm:"foreignKey:BagID" json:"products,omitempty"`
	StoreTransferBags []StoreTransferBag `gorm:"foreignKey:BagID" json:"store_transfer_bags,omitempty"`
}

// BeforeCreate generates UUID before creating
func (b *Bag) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
