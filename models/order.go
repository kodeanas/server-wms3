package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID         uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code       string          `gorm:"size:255;not null;index"`
	Type       string          `gorm:"type:varchar(50);check:type IN ('regular','cargo','qcd','scrap')"`
	BuyerID    uuid.UUID       `gorm:"type:uuid;index"`
	ClassID    uuid.UUID       `gorm:"type:uuid;index"`
	UserID     *uuid.UUID      `gorm:"type:uuid;index"` // nullable
	Status     string          `gorm:"type:varchar(50);check:status IN ('progress','pending','done','cancel');index"`
	IsTax      bool            `gorm:"default:false"`
	Tax        float64         `gorm:"type:decimal(15,2)"`
	TaxValue   float64         `gorm:"type:decimal(15,2)"`
	TotalPrice float64         `gorm:"type:decimal(15,2)"`
	TotalBox   int             `gorm:"type:int"`
	PriceBox   float64         `gorm:"type:decimal(15,2)"`
	GrandTotal float64         `gorm:"type:decimal(15,2)"`
	CreatedAt  time.Time       `gorm:"autoCreateTime"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt  `gorm:"index"`
}
