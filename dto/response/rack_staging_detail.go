package dto

type RackStagingDetailResponse struct {
	Code                string  `json:"code"`
	RackDisplayName     string  `json:"rack_display_name"`
	CreatedAt           string  `json:"created_at"`
	IsMoved             bool    `json:"is_moved"`
	TotalItem           int     `json:"total_item"`
	TotalPriceWarehouse float64 `json:"total_price_warehouse"`
}
