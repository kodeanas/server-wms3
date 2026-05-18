package dto

// BulkSummaryAllResponse is summary all for BULK documents and product info
type BulkSummaryAllResponse struct {
	TotalDocumentUpload int     `json:"total_document_upload"`
	TotalProductMasuk   int     `json:"total_product_masuk"`
	TotalHargaMasuk     float64 `json:"total_harga_masuk"`
}
