package dto

type ListCategory struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Discount  *int     `json:"discount"`
	MinPrice  *float64 `json:"min_price"`
	MaxPrice  *float64 `json:"max_price"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	DeletedAt string   `json:"deleted_at"`
}

type ListCategorySelect struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Discount *int     `json:"discount"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Status   string   `json:"status"`
}
