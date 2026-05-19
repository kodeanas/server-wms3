package dto

import "time"

// BulkSummaryAllResponse is summary all for BULK documents and product info
type BulkSummaryAllResponse struct {
	TotalDocumentUpload int     `json:"total_document_upload"`
	TotalProductMasuk   int     `json:"total_product_masuk"`
	TotalHargaMasuk     float64 `json:"total_harga_masuk"`
}

// BulkDocumentDetailResponse contains core detail info for one BULK document.
type BulkDocumentDetailResponse struct {
	ID         string  `json:"id"`
	Code       string  `json:"code"`
	Nama       string  `json:"nama"`
	TotalPrice float64 `json:"total_price"`
	TotalItem  int     `json:"total_item"`
}

// BulkDocumentSummaryItemResponse is grouped summary inside a BULK document detail.
type BulkDocumentSummaryItemResponse struct {
	Label          string  `json:"label"`
	Item           int     `json:"item"`
	Price          float64 `json:"price"`
	PriceWarehouse float64 `json:"price_warehouse"`
}

// BulkProductDocumentItemResponse is one product row inside BULK document.
type BulkProductDocumentItemResponse struct {
	ID          string     `json:"id"`
	Barcode     string     `json:"barcode"`
	Name        string     `json:"name"`
	Item        int        `json:"item"`
	Price       float64    `json:"price"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
	DateScanned *time.Time `json:"date_scanned"`
}
