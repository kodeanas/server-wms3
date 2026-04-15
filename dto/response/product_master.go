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

type ProductMasterDetailResponse struct {
	ID               string                        `json:"id"`
	DocumentID       string                        `json:"document_id"`
	DocumentCode     string                        `json:"document_code"`
	DocumentName     string                        `json:"document_name"`
	Barcode          string                        `json:"barcode"`
	BarcodeWarehouse string                        `json:"barcode_warehouse"`
	Name             string                        `json:"name"`
	NameWarehouse    string                        `json:"name_warehouse"`
	Item             int                           `json:"item"`
	ItemWarehouse    int                           `json:"item_warehouse"`
	Price            float64                       `json:"price"`
	PriceWarehouse   float64                       `json:"price_warehouse"`
	CategoryID       *string                       `json:"category_id"`
	StickerID        *string                       `json:"sticker_id"`
	TypeOut          string                        `json:"type_out"`
	Location         string                        `json:"location"`
	CreatedAt        time.Time                     `json:"created_at"`
	Additional       ProductMasterDetailAdditional `json:"additional"`
}

type ProductMasterDetailAdditional struct {
	Document *ProductMasterDocumentAdditional `json:"document"`
	Sticker  *ProductMasterStickerAdditional  `json:"sticker"`
	Category *ProductMasterCategoryAdditional `json:"category"`
}

type ProductMasterDocumentAdditional struct {
	Code     string `json:"code"`
	NameFile string `json:"name_file"`
}

type ProductMasterStickerAdditional struct {
	Name    string `json:"name"`
	CodeHex string `json:"code_hex"`
	Type    string `json:"type"`
}

type ProductMasterCategoryAdditional struct {
	Name    string `json:"name"`
	Diskon  int    `json:"diskon"`
	Type    string `json:"type"`
	CodeHex string `json:"code_hex"`
}
