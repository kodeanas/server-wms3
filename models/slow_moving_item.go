package models

import (
	"time"

	"github.com/google/uuid"
)

type SlowMovingItem struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SlowMovingID    string    `gorm:"type:char(36)" json:"slow_moving_id"`
	ProductMasterID string    `gorm:"type:char(36)" json:"product_master_id"`
	IsDamaged       bool      `gorm:"default:false" json:"is_damaged"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
