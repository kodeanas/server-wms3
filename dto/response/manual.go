package dto

type ProductManualResponse struct {
	ID         string  `json:"id"`
	DocumentID string  `json:"document_id"`
	Barcode    string  `json:"barcode"`
	Name       string  `json:"name"`
	Item       int     `json:"item"`
	Price      float64 `json:"price"`
	Status     string  `json:"status"`
	Note       string  `json:"note"`
	CreatedAt  string  `json:"created_at"`

	PriceWarehouse *float64 `json:"price_warehouse"`
	CategoryName   *string  `json:"category_name"`
	StickerName    *string  `json:"sticker_name"`
}
