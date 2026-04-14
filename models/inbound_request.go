package models

type InboundRequest struct {
	Name       string  `json:"name" binding:"required"`
	Item       int     `json:"item" binding:"required,gt=0"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID *string `json:"category_id,omitempty"`
	StickerID  *string `json:"sticker_id,omitempty"`
	Status     string  `json:"status" binding:"required,oneof=good abnormal damaged non"`
	Note       *string `json:"note,omitempty"`
}
