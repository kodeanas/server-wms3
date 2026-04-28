package dto

type BuyerClassInfoResponse struct {
	BuyerID       string  `json:"buyer_id"`
	BuyerName     string  `json:"buyer_name"`
	ClassName     string  `json:"class_name"`
	ClassDiscount float64 `json:"class_discount"`
	NextClassNote string  `json:"next_class_note"`
}

type ProductOrderDetail struct {
	ProductOrderID string  `json:"product_order_id"`
	ProductID      string  `json:"product_id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Discount       float64 `json:"discount"`
	Qty            int     `json:"qty"`
	Subtotal       float64 `json:"subtotal"`
}

type DiscountOrderDetail struct {
	DiscountOrderID string  `json:"discount_order_id"`
	Type            string  `json:"type"`
	Name            string  `json:"name"`
	IsNominal       bool    `json:"is_nominal"`
	ValueNominal    float64 `json:"value_nominal"`
	ValuePercent    int     `json:"value_percentage"`
	Usage           float64 `json:"usage"`
}

type OutboundRegulerOrderDetail struct {
	OrderID        string                `json:"order_id"`
	Code           string                `json:"code"`
	BuyerID        string                `json:"buyer_id"`
	BuyerName      string                `json:"buyer_name"`
	Status         string                `json:"status"`
	IsTax          bool                  `json:"is_tax"`
	Tax            float64               `json:"tax"`
	TaxValue       float64               `json:"tax_value"`
	TotalBox       int                   `json:"total_box"`
	PriceBox       float64               `json:"price_box"`
	ProductOrders  []ProductOrderDetail  `json:"product_orders"`
	DiscountOrders []DiscountOrderDetail `json:"discount_orders"`
	TotalPrice     float64               `json:"total_price"`
	GrandTotal     float64               `json:"grand_total"`
	CreatedAt      string                `json:"created_at"`
}
