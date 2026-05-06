package dto

// BuyerResponse for custom API response
type BuyerResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Email     string               `json:"email"`
	Phone     string               `json:"phone"`
	ClassID   string               `json:"class_id"`
	Address   string               `json:"address"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
	DeletedAt interface{}          `json:"deleted_at"`
	Class     *ClassSimpleResponse `json:"class,omitempty"`
}

type ClassSimpleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
