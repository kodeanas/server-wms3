package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Rack represents a warehouse rack
type Rack struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // enum: display, staging
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	RackID    string    `gorm:"type:varchar(255)" json:"rack_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Products []ProductMaster `gorm:"foreignKey:RackID" json:"products,omitempty"`
}

// BeforeCreate generates UUID before creating
func (r *Rack) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// Store represents a warehouse store
type Store struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20);not null;index" json:"phone"`
	Email     sql.NullString `gorm:"type:varchar(255)" json:"email"`
	Address   sql.NullString `gorm:"type:text" json:"address"`
	UserID    string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	IsCashier bool           `json:"is_cashier"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	User           *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	StoreCrew      []StoreCrew        `gorm:"foreignKey:StoreID" json:"store_crew,omitempty"`
	StoreTransfers []StoreTransfer    `gorm:"foreignKey:StoreID" json:"store_transfers,omitempty"`
	TransferBags   []StoreTransferBag `gorm:"foreignKey:StoreID" json:"transfer_bags,omitempty"`
}

// BeforeCreate generates UUID before creating
func (s *Store) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// StoreCrew represents store staff
type StoreCrew struct {
	ID        string         `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20);not null;index" json:"phone"`
	Email     sql.NullString `gorm:"type:varchar(255)" json:"email"`
	Address   sql.NullString `gorm:"type:text" json:"address"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	StoreID   string         `gorm:"type:uuid;not null;index" json:"store_id"`
	IsCashier bool           `json:"is_cashier"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

// BeforeCreate generates UUID before creating
func (sc *StoreCrew) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == "" {
		sc.ID = uuid.New().String()
	}
	return nil
}

// StoreTransfer represents warehouse transfer
type StoreTransfer struct {
	ID                  string         `gorm:"primaryKey;type:uuid" json:"id"`
	StoreID             string         `gorm:"type:uuid;not null;index" json:"store_id"`
	DateTime            time.Time      `json:"date_time"`
	TotalItem           sql.NullInt64  `json:"total_item"`
	TotalPrice          sql.NullString `gorm:"type:numeric(19,2)" json:"total_price"`
	TotalPriceWarehouse sql.NullString `gorm:"type:numeric(19,2)" json:"total_price_warehouse"`
	Status              string         `gorm:"type:varchar(50);default:'on_progress'" json:"status"` // enum
	UserID              string         `gorm:"type:uuid;index" json:"user_id"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`

	// Relationships
	Store             *Store             `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	User              *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	StoreTransferBags []StoreTransferBag `gorm:"foreignKey:StoreTransferID" json:"transfer_bags,omitempty"`
}

// BeforeCreate generates UUID before creating
func (st *StoreTransfer) BeforeCreate(tx *gorm.DB) error {
	if st.ID == "" {
		st.ID = uuid.New().String()
	}
	return nil
}

// StoreTransferBag represents items in a transfer
type StoreTransferBag struct {
	ID              string         `gorm:"primaryKey;type:uuid" json:"id"`
	StoreTransferID string         `gorm:"type:uuid;not null;index" json:"store_transfer_id"`
	BagID           string         `gorm:"type:uuid;index" json:"bag_id"`
	Quantity        int            `json:"quantity"`
	TotalPrice      sql.NullString `gorm:"type:numeric(19,2)" json:"total_price"`
	TotalCogs       sql.NullString `gorm:"type:numeric(19,2)" json:"total_cogs"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relationships
	StoreTransfer *StoreTransfer `gorm:"foreignKey:StoreTransferID" json:"store_transfer,omitempty"`
	Bag           *Bag           `gorm:"foreignKey:BagID" json:"bag,omitempty"`
	Store         *Store         `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

// BeforeCreate generates UUID before creating
func (stb *StoreTransferBag) BeforeCreate(tx *gorm.DB) error {
	if stb.ID == "" {
		stb.ID = uuid.New().String()
	}
	return nil
}
