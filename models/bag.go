package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bag struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code            string         `gorm:"size:255;not null" json:"code"`
	Type            string         `gorm:"size:50" json:"type"`
	UserID          *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	IsMoved         bool           `gorm:"default:false" json:"is_moved"`
	DateOut         *time.Time     `json:"date_out"`
	CargoID         *uuid.UUID     `gorm:"type:uuid" json:"cargo_id"`
	TransferStoreID *uuid.UUID     `gorm:"type:uuid" json:"transfer_store_id"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
