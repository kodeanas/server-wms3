package dto

// ClassResponse for custom API response
type ClassResponse struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	MinOrder            int         `json:"min_order"`
	Disc                int         `json:"disc"`
	MinTransactionValue string      `json:"min_transaction_value"`
	Week                int         `json:"week"`
	Iteration           int         `json:"iteration"`
	Status              string      `json:"status"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	DeletedAt           interface{} `json:"deleted_at"`
}
