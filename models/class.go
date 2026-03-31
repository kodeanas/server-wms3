package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Class represents a buyer class/tier
type Class struct {
	ID                  string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name                string    `gorm:"type:varchar(255);not null" json:"name"`
	MinOrder            int       `gorm:"not null" json:"min_order"`
	Discount            int       `json:"discount"`
	MinTransactionValue string    `gorm:"type:numeric(19,2)" json:"min_transaction_value"`
	Week                int       `json:"week"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relationships
	Buyers        []Buyer        `gorm:"foreignKey:ClassID" json:"buyers,omitempty"`
	UserClassLogs []UserClassLog `gorm:"foreignKey:NewClassID" json:"user_class_logs,omitempty"`
}

// BeforeCreate generates UUID before creating
func (cl *Class) BeforeCreate(tx *gorm.DB) error {
	if cl.ID == "" {
		cl.ID = uuid.New().String()
	}
	return nil
}

// UserClassLog represents changes in user class
type UserClassLog struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	PrevClassID string    `gorm:"type:uuid" json:"prev_class_id"`
	NewClassID  string    `gorm:"type:uuid;not null;index" json:"new_class_id"`
	OrderID     string    `gorm:"type:uuid;not null;index" json:"order_id"`
	ChangeType  string    `gorm:"type:varchar(50);not null" json:"change_type"` // enum: upgrade, downgrade
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	User      *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PrevClass *Class `gorm:"foreignKey:PrevClassID" json:"prev_class,omitempty"`
	NewClass  *Class `gorm:"foreignKey:NewClassID" json:"new_class,omitempty"`
	Order     *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// BeforeCreate generates UUID before creating
func (ucl *UserClassLog) BeforeCreate(tx *gorm.DB) error {
	if ucl.ID == "" {
		ucl.ID = uuid.New().String()
	}
	return nil
}
