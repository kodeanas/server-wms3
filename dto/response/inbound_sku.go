package dto

type InboundSKUDocumentDTO struct {
	ID        string  `json:"id"`
	FileName  string  `json:"file_name"`
	FileItem  int     `json:"file_item"`
	FilePrice int     `json:"file_price"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
	UserID    *string `json:"user_id"`
	Supplier  string  `json:"supplier"`
}

type ProductPendingDTO struct {
	ID          string  `json:"id"`
	Barcode     string  `json:"barcode"`
	Name        string  `json:"name"`
	Item        int     `json:"item"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	Note        string  `json:"note"`
	DateScanned *string `json:"date_scanned"`
	ItemGood    int     `json:"item_good"`
	ItemDamaged int     `json:"item_damaged"`
}

// SkuSummaryAllResponse is summary all for SKU documents and product info
type SkuSummaryAllResponse struct {
	TotalDocumentUpload int     `json:"total_document_upload"`
	TotalProductMasuk   int     `json:"total_product_masuk"`
	TotalHargaMasuk     float64 `json:"total_harga_masuk"`
}
