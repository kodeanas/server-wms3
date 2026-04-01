package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductLog struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductMasterID string    `gorm:"type:char(36)" json:"product_master_id"`
	DataName        string    `gorm:"size:255" json:"data_name"`
	PrevData        string    `gorm:"type:text" json:"prev_data"`
	NewData         string    `gorm:"type:text" json:"new_data"`
	UserID          string    `gorm:"type:char(36)" json:"user_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
