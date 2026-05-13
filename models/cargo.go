package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cargo struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code      string         `gorm:"size:255;not null;uniqueIndex" json:"code"`
	Status    string         `gorm:"size:50" json:"status"`
	IsSale    bool           `gorm:"default:false" json:"is_sale"`
	IsOnline  bool           `gorm:"default:false" json:"is_online"`
	UserID    *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
