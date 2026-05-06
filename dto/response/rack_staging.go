package dto

type CreateRackStagingRequest struct {
	RackDisplayID string `json:"rack_display_id" binding:"required"`
}

type RackStagingDetailResponse struct {
	Code                string  `json:"code"`
	RackDisplayName     string  `json:"rack_display_name"`
	CreatedAt           string  `json:"created_at"`
	IsMoved             bool    `json:"is_moved"`
	TotalItem           int     `json:"total_item"`
	TotalPriceWarehouse float64 `json:"total_price_warehouse"`
}

type RackStagingResponse struct {
	ID            string `json:"id"`
	RackDisplayID string `json:"rack_display_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	IsMoved       bool   `json:"is_moved"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
