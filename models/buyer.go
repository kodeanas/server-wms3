package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Buyer struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"size:255" json:"name"`
	Email     string         `gorm:"size:255" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	ClassID   string         `gorm:"type:char(36)" json:"class_id"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// CreateBuyerPayload request payload.
type CreateBuyerPayload struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	ClassID string `json:"class_id"`
	Address string `json:"address"`
}

// UpdateBuyerPayload request payload for update.
type UpdateBuyerPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	ClassID string `json:"class_id"`
	Address string `json:"address"`
}
