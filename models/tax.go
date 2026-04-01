package models

import (
	"time"

	"github.com/google/uuid"
)

type Tax struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Tax       int       `gorm:"not null" json:"tax"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
