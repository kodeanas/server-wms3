package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Discount  int       `gorm:"default:0" json:"discount"`
	MinPrice  string    `gorm:"type:numeric(19,2)" json:"min_price"`
	MaxPrice  string    `gorm:"type:numeric(19,2)" json:"max_price"`
	Status    string    `gorm:"type:varchar(50);default:'active'" json:"status"` // enum: active, inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Products []ProductMaster `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	Stickers []Sticker       `gorm:"many2many:category_stickers" json:"stickers,omitempty"`
}

// BeforeCreate generates UUID before creating
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Sticker represents a product sticker/label
type Sticker struct {
	ID         string    `gorm:"primaryKey;type:uuid" json:"id"`
	CodeHex    string    `gorm:"type:varchar(255);not null" json:"code_hex"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug       string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Type       string    `gorm:"type:varchar(50)" json:"type"` // enum: big, small
	FixedPrice int       `json:"fixed_price"`
	MinPrice   string    `gorm:"type:numeric(19,2)" json:"min_price"`
	MaxPrice   string    `gorm:"type:numeric(19,2)" json:"max_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	ProductMasters []ProductMaster `gorm:"many2many:product_stickers" json:"product_masters,omitempty"`
	Categories     []Category      `gorm:"many2many:category_stickers" json:"categories,omitempty"`
}

// BeforeCreate generates UUID before creating
func (s *Sticker) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}
