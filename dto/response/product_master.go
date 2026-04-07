package dto

import "time"

type ProductMasterRegulerResponse struct {
	ID               string    `json:"id"`
	DocumentID       string    `json:"document_id"`
	Barcode          string    `json:"barcode"`
	BarcodeWarehouse string    `json:"barcode_warehouse"`
	Name             string    `json:"name"`
	NameWarehouse    string    `json:"name_warehouse"`
	Item             int       `json:"item"`
	ItemWarehouse    int       `json:"item_warehouse"`
	Price            float64   `json:"price"`
	PriceWarehouse   float64   `json:"price_warehouse"`
	CategoryID       *string   `json:"category_id"`
	ProductPendingID *string   `json:"product_pending_id"`
	IsSKU            bool      `json:"is_sku"`
	Location         string    `json:"location"`
	TypeOut          string    `json:"type_out"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CategoryName     *string   `json:"category_name"`
}

type ProductMasterStickerResponse struct {
	ID               string    `json:"id"`
	DocumentID       string    `json:"document_id"`
	Barcode          string    `json:"barcode"`
	BarcodeWarehouse string    `json:"barcode_warehouse"`
	Name             string    `json:"name"`
	NameWarehouse    string    `json:"name_warehouse"`
	Item             int       `json:"item"`
	ItemWarehouse    int       `json:"item_warehouse"`
	Price            float64   `json:"price"`
	PriceWarehouse   float64   `json:"price_warehouse"`
	StickerID        *string   `json:"sticker_id"`
	ProductPendingID *string   `json:"product_pending_id"`
	IsSKU            bool      `json:"is_sku"`
	Location         string    `json:"location"`
	TypeOut          string    `json:"type_out"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	StickerName      *string   `json:"sticker_name"`
	StickerCodeHex   *string   `json:"sticker_code_hex"`
}
