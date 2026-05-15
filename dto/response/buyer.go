package dto

// BuyerResponse for custom API response
type BuyerResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	ClassName string      `json:"class_name"`
	Address   string      `json:"address"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	DeletedAt interface{} `json:"deleted_at"`
}

type BuyerClassResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	ClassName string `json:"class_name"`
	Address   string `json:"address"`
}
