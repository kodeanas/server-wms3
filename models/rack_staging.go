package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RackStaging struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RackDisplayID uuid.UUID      `gorm:"type:uuid;not null" json:"rack_display_id"`
	Code          string         `gorm:"size:255;not null;unique" json:"code"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	IsMoved       bool           `gorm:"default:false" json:"is_moved"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
