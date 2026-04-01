package models

import (
	"time"

	"github.com/google/uuid"
)

type Buyer struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Nama      string    `gorm:"size:255" json:"nama"`
	Email     string    `gorm:"size:255" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	ClassID   string    `gorm:"type:char(36)" json:"class_id"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
