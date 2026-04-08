package dto

type InboundBastSummaryResponse struct {
	TotalFileUpload       int     `json:"total_file_upload"`
	TotalFileMasihProses  int     `json:"total_file_masih_proses"`
	TotalItemTerScan      int     `json:"total_item_ter_scan"`
	TotalHargaAsalTerscan float64 `json:"total_harga_asal_terscan"`
}
