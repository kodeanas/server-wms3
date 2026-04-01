package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;unique;not null" json:"email"`
	Phone     string    `gorm:"size:20;unique;not null" json:"phone"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
