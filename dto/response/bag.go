package dto

import "time"

type WholesaleBagListResponse struct {
	ID        string      `json:"id"`
	Code      string      `json:"code"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt interface{} `json:"deleted_at"`
}
