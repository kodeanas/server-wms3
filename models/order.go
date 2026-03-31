package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents a customer order
type Order struct {
	ID                  string         `gorm:"primaryKey;type:uuid" json:"id"`
	Code                string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"code"`
	UserID              string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Type                string         `gorm:"type:varchar(50)" json:"type"` // enum: regular, wcargo, wjcd, wacarg
	BuyerID             string         `gorm:"type:uuid;not null;index" json:"buyer_id"`
	DateTime            time.Time      `json:"date_time"`
	TotalItem           int            `json:"total_item"`
	TotalPriceProduct   sql.NullString `gorm:"type:numeric(19,2)" json:"total_price_product"`
	TotalPPN            sql.NullString `gorm:"type:numeric(19,2)" json:"total_ppn"`
	CartonBoxPrice      sql.NullString `gorm:"type:numeric(19,2)" json:"carton_box_price"`
	QuantityCartonBox   int            `json:"quantity_carton_box"`
	Voucher             sql.NullString `gorm:"type:numeric(19,2)" json:"voucher"`
	FixedDiscount       int            `json:"fixed_discount"`
	ClassDiscount       int            `json:"class_discount"`
	ClassDiscountAmount sql.NullString `gorm:"type:numeric(19,2)" json:"class_discount_amount"`
	GrandTotal          sql.NullString `gorm:"type:numeric(19,2)" json:"grand_total"`
	Status              string         `gorm:"type:varchar(50);default:'on_progress'" json:"status"` // enum
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`

	// Relationships
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Buyer     *Buyer         `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Items     []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Cargos    []OrderCargo   `gorm:"foreignKey:OrderID" json:"cargos,omitempty"`
	ClassLogs []UserClassLog `gorm:"foreignKey:OrderID" json:"class_logs,omitempty"`
}

// BeforeCreate generates UUID before creating
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// OrderItem represents items in an order
type OrderItem struct {
	ID              string         `gorm:"primaryKey;type:uuid" json:"id"`
	ProductMasterID string         `gorm:"type:uuid;not null;index" json:"product_master_id"`
	OrderID         string         `gorm:"type:uuid;not null;index" json:"order_id"`
	Price           sql.NullString `gorm:"type:numeric(19,2)" json:"price"`
	PriceWarehouse  sql.NullString `gorm:"type:numeric(19,2)" json:"price_warehouse"`
	PriceCut        sql.NullString `gorm:"type:numeric(19,2)" json:"price_cut"`
	PriceFinal      sql.NullString `gorm:"type:numeric(19,2)" json:"price_final"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// Relationships
	Order         *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductMaster *ProductMaster `gorm:"foreignKey:ProductMasterID" json:"product_master,omitempty"`
	OrderCargos   []OrderCargo   `gorm:"foreignKey:OrderItemID" json:"order_cargos,omitempty"`
}

// BeforeCreate generates UUID before creating
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = uuid.New().String()
	}
	return nil
}

// OrderCargo represents cargo in an order
type OrderCargo struct {
	ID             string         `gorm:"primaryKey;type:uuid" json:"id"`
	CargoID        string         `gorm:"type:uuid;not null;index" json:"cargo_id"`
	OrderID        string         `gorm:"type:uuid;not null;index" json:"order_id"`
	OrderItemID    string         `gorm:"type:uuid;not null;index" json:"order_item_id"`
	Price          sql.NullString `gorm:"type:numeric(19,2)" json:"price"`
	PriceWarehouse sql.NullString `gorm:"type:numeric(19,2)" json:"price_warehouse"`
	PriceCut       sql.NullString `gorm:"type:numeric(19,2)" json:"price_cut"`
	PriceFinal     sql.NullString `gorm:"type:numeric(19,2)" json:"price_final"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`

	// Relationships
	Cargo     *Cargo     `gorm:"foreignKey:CargoID" json:"cargo,omitempty"`
	Order     *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	OrderItem *OrderItem `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
}

// BeforeCreate generates UUID before creating
func (oc *OrderCargo) BeforeCreate(tx *gorm.DB) error {
	if oc.ID == "" {
		oc.ID = uuid.New().String()
	}
	return nil
}
