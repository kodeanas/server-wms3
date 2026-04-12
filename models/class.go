package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Class struct {
	ID                  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name                string         `gorm:"size:255;not null" json:"name"`
	MinOrder            int            `gorm:"not null" json:"min_order"`
	Disc                int            `gorm:"default:0" json:"disc"`
	MinTransactionValue float64        `gorm:"type:decimal(15,2)" json:"min_transaction_value"`
	Week                int            `json:"week"`
	Iteration           int            `json:"iteration"`
	Status              string         `gorm:"size:50;check:status IN ('active', 'inactive');default:'active'" json:"status"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
