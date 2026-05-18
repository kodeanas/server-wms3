package dto

import "time"

type BastDocumentResponse struct {
	ID        string      `json:"id"`
	Code      string      `json:"code"`
	FileName  string      `json:"file_name"`
	FileItem  int         `json:"file_item"`
	FilePrice int         `json:"file_price"`
	Status    string      `json:"status"`
	UserID    *string     `json:"user_id"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	DeletedAt interface{} `json:"deleted_at"`
	DateStop  interface{} `json:"date_stop"`
}

type BastProductResponse struct {
	ID          string     `json:"id"`
	Barcode     string     `json:"barcode"`
	Name        string     `json:"name"`
	Item        int        `json:"item"`
	Price       float64    `json:"price"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
	DateScanned *time.Time `json:"date_scanned"`
}

type InboundBastSummaryResponse struct {
	TotalFileUpload       int     `json:"total_file_upload"`
	TotalFileMasihProses  int     `json:"total_file_masih_proses"`
	TotalItemTerScan      int     `json:"total_item_ter_scan"`
	TotalHargaAsalTerscan float64 `json:"total_harga_asal_terscan"`
}

// BastSummaryAllResponse is summary all for BAST documents and product status breakdown
type BastSummaryAllResponse struct {
	TotalDocumentInbound int `json:"total_document_inbound"`
	TotalDocumentScanned int `json:"total_document_scanned"`
	TotalProductGood     int `json:"total_product_good"`
	TotalProductDamaged  int `json:"total_product_damaged"`
	TotalProductAbnormal int `json:"total_product_abnormal"`
	TotalProductNon      int `json:"total_product_non"`
}
