package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Taxes         []Tax           `gorm:"foreignKey:UserID" json:"taxes,omitempty"`
	UserClassLogs []UserClassLog  `gorm:"foreignKey:UserID" json:"user_class_logs,omitempty"`
	Products      []ProductMaster `gorm:"foreignKey:UserID" json:"products,omitempty"`
	Store         *Store          `gorm:"foreignKey:UserID" json:"store,omitempty"`
	Cargos        []Cargo         `gorm:"foreignKey:UserID" json:"cargos,omitempty"`
	Orders        []Order         `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	SlowMovings   []SlowMoving    `gorm:"foreignKey:UserID" json:"slow_movings,omitempty"`
}

// BeforeCreate generates UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Tax represents a tax rate in the system
type Tax struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Tax       int       `gorm:"not null" json:"tax"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate generates UUID before creating
func (t *Tax) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// Buyer represents a buyer profile
type Buyer struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     sql.NullString `gorm:"type:varchar(255)" json:"email"`
	Phone     string         `gorm:"type:varchar(20);not null" json:"phone"`
	ClassID   string         `gorm:"type:uuid;index" json:"class_id"`
	Address   sql.NullString `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	Class  *Class  `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Orders []Order `gorm:"foreignKey:BuyerID" json:"orders,omitempty"`
}

// BeforeCreate generates UUID before creating
func (b *Buyer) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
