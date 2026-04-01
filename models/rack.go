package models

import (
	"time"

	"github.com/google/uuid"
)

type Rack struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type      string    `gorm:"size:50" json:"type"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	RackID    string    `gorm:"size:255;not null" json:"rack_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
