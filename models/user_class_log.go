package models

import (
	"time"

	"github.com/google/uuid"
)

type UserClassLog struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BuyerID     string    `gorm:"type:char(36)" json:"buyer_id"`
	PrevClassID string    `gorm:"type:char(36)" json:"prev_class_id"`
	NewClassID  string    `gorm:"type:char(36)" json:"new_class_id"`
	OrderID     string    `gorm:"type:char(36)" json:"order_id"`
	ChangeType  string    `gorm:"size:255" json:"change_type"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
