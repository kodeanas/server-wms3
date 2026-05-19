package dto

import "time"

// RackDisplayDetailResponse is the response for rack display detail
// with total_item, total_price, total_price_warehouse
//
type RackDisplayDetailResponse struct {
	ID                  string                            `json:"id"`
	Code                string                            `json:"code"`
	Name                string                            `json:"name"`
	CreatedAt           time.Time                         `json:"created_at"`
	TotalItem           int                               `json:"total_item"`
	TotalPrice          float64                           `json:"total_price"`
	TotalPriceWarehouse float64                           `json:"total_price_warehouse"`
	Summary             []RackDisplaySummaryItemResponse  `json:"summary"`
}

// RackDisplaySummaryItemResponse is grouped summary inside a rack display detail.
type RackDisplaySummaryItemResponse struct {
	Label          string  `json:"label"`
	Item           int     `json:"item"`
	Price          float64 `json:"price"`
	PriceWarehouse float64 `json:"price_warehouse"`
}

// RackDisplaySummaryAllResponse is summary across all rack display product lists.
type RackDisplaySummaryAllResponse struct {
	TotalRack         int64   `json:"total_rack"`
	TotalItem         int64   `json:"total_item"`
	TotalPrice        float64 `json:"total_price"`
	TotalItemAndPrice float64 `json:"total_item_and_price"`
}

// RackProductSummaryResponse is summary of products in a rack display (virtual rack).
// Contains total_item, total_price, total_price_warehouse, and total_item_and_price.
type RackProductSummaryResponse struct {
	TotalItem           int64   `json:"total_item"`
	TotalPrice          float64 `json:"total_price"`
	TotalPriceWarehouse float64 `json:"total_price_warehouse"`
	TotalItemAndPrice   float64 `json:"total_item_and_price"`
}
