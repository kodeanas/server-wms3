package dto

import "time"

// RackDisplayDetailResponse is the response for rack display detail
// with total_item, total_price, total_price_warehouse
//
type RackDisplayDetailResponse struct {
	ID                  string    `json:"id"`
	Code                string    `json:"code"`
	Name                string    `json:"name"`
	CreatedAt           time.Time `json:"created_at"`
	TotalItem           int       `json:"total_item"`
	TotalPrice          float64   `json:"total_price"`
	TotalPriceWarehouse float64   `json:"total_price_warehouse"`
}
