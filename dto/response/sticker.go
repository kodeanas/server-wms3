package dto

// StickerResponse adalah response untuk detail/single sticker
type StickerResponse struct {
	ID         string   `json:"id"`
	CodeHex    string   `json:"code_hex"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Status     string   `json:"status"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	DeletedAt  string   `json:"deleted_at"`
}

// ListStickerSelect adalah response untuk list sticker (select/dropdown)
type ListStickerSelect struct {
	ID         string   `json:"id"`
	CodeHexx   string   `json:"code_hex"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	Type       string   `json:"type"`
	FixedPrice *int     `json:"fixed_price"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Status     string   `json:"status"`
}
