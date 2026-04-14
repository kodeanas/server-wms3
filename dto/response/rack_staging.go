package dto

type CreateRackStagingRequest struct {
	RackDisplayID string `json:"rack_display_id" binding:"required"`
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
