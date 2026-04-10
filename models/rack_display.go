package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RackDisplay struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code      string         `gorm:"size:255;not null;unique" json:"code"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
