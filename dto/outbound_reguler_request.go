package dto

type ScanProductRequest struct {
	Barcode string `json:"barcode" binding:"required"`
	OrderID string `json:"order_id"`
	BuyerID string `json:"buyer_id"`
}

type AddProductRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	ProductID string `json:"product_id" binding:"required"`
	Qty       int    `json:"qty" binding:"required"`
}

type AddDiscountRequest struct {
	OrderID      string  `json:"order_id" binding:"required"`
	Type         string  `json:"type" binding:"required"` // voucher, rank, additional
	Name         string  `json:"name"`
	IsNominal    bool    `json:"is_nominal"`
	ValueNominal float64 `json:"value_nominal"`
	ValuePercent int     `json:"value_percentage"`
}

type UpdateTaxRequest struct {
	OrderID  string  `json:"order_id" binding:"required"`
	IsTax    bool    `json:"is_tax"`
	Tax      float64 `json:"tax"`
	TaxValue float64 `json:"tax_value"`
}

type UpdateBoxRequest struct {
	OrderID  string  `json:"order_id" binding:"required"`
	TotalBox int     `json:"total_box"`
	PriceBox float64 `json:"price_box"`
}
