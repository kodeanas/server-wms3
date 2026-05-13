package dto

import "time"

type CargoListResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	IsSale    bool      `json:"is_sale"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BagItemResponse struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Type    string `json:"type"`
	IsMoved bool   `json:"is_moved"`
}

type CargoDetailResponse struct {
	ID         string            `json:"id"`
	Code       string            `json:"code"`
	Status     string            `json:"status"`
	IsSale     bool              `json:"is_sale"`
	IsOnline   bool              `json:"is_online"`
	CreatedAt  string            `json:"created_at"`
	TotalBag   int               `json:"total_bag"`
	TotalItem  int               `json:"total_item"`
	TotalPrice float64           `json:"total_price"`
	Bags       []BagItemResponse `json:"bags"`
}
