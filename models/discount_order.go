package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountOrder struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID         string         `gorm:"type:char(36)" json:"order_id"`
	Type            string         `gorm:"size:50" json:"type"`
	Name            string         `gorm:"size:255" json:"name"`
	IsNominal       bool           `gorm:"default:true" json:"is_nominal"`
	ValueNominal    float64        `gorm:"type:decimal(15,2)" json:"value_nominal"`
	ValuePercentage int            `json:"value_percentage"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
