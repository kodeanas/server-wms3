package repositories

import (
	"time"

	"gorm.io/gorm"
)

type InboundBastSummaryFilter struct {
	Date      *time.Time
	DateStart *time.Time
	DateEnd   *time.Time
}

type InboundBastSummaryResult struct {
	TotalFileUpload       int
	TotalFileMasihProses  int
	TotalItemTerScan      int
	TotalHargaAsalTerscan float64
}

func GetInboundBastSummary(db *gorm.DB, filter InboundBastSummaryFilter) (InboundBastSummaryResult, error) {
	var result InboundBastSummaryResult

	// Query for TotalFileUpload
	fileUploadQuery := db.Table("product_documents")
	var totalFileUpload int64
	if filter.Date != nil {
		fileUploadQuery = fileUploadQuery.Where("DATE(created_at) = ?", filter.Date.Format("2006-01-02"))
	} else if filter.DateStart != nil && filter.DateEnd != nil {
		fileUploadQuery = fileUploadQuery.Where("created_at > ? AND created_at < ?", filter.DateStart.Format("2006-01-02"), filter.DateEnd.Format("2006-01-02"))
	} else {
		fileUploadQuery = fileUploadQuery.Where("DATE(created_at) = ?", time.Now().Format("2006-01-02"))
	}
	fileUploadQuery.Count(&totalFileUpload)
	result.TotalFileUpload = int(totalFileUpload)

	// Query for TotalFileMasihProses
	fileProsesQuery := db.Table("product_documents")
	var totalFileMasihProses int64
	if filter.Date != nil {
		fileProsesQuery = fileProsesQuery.Where("created_at < ? AND date_out > ?", filter.Date.Format("2006-01-02"), filter.Date.Format("2006-01-02"))
	} else if filter.DateStart != nil && filter.DateEnd != nil {
		fileProsesQuery = fileProsesQuery.Where("created_at < ? AND date_out > ?", filter.DateStart.Format("2006-01-02"), filter.DateEnd.Format("2006-01-02"))
	} else {
		fileProsesQuery = fileProsesQuery.Where("status = ?", "progress")
	}
	fileProsesQuery.Count(&totalFileMasihProses)
	result.TotalFileMasihProses = int(totalFileMasihProses)

	// Query for TotalItemTerScan
	itemScanQuery := db.Table("product_pendings")
	var totalItemTerScan int64
	if filter.Date != nil {
		itemScanQuery = itemScanQuery.Where("DATE(date_scanned) = ?", filter.Date.Format("2006-01-02"))
	} else if filter.DateStart != nil && filter.DateEnd != nil {
		itemScanQuery = itemScanQuery.Where("date_scanned > ? AND date_scanned < ?", filter.DateStart.Format("2006-01-02"), filter.DateEnd.Format("2006-01-02"))
	} else {
		itemScanQuery = itemScanQuery.Where("DATE(date_scanned) = ?", time.Now().Format("2006-01-02"))
	}
	itemScanQuery.Count(&totalItemTerScan)
	result.TotalItemTerScan = int(totalItemTerScan)

	// Query for TotalHargaAsalTerscan
	hargaQuery := db.Table("product_pendings")
	if filter.Date != nil {
		hargaQuery = hargaQuery.Where("DATE(date_scanned) = ?", filter.Date.Format("2006-01-02"))
	} else if filter.DateStart != nil && filter.DateEnd != nil {
		hargaQuery = hargaQuery.Where("date_scanned > ? AND date_scanned < ?", filter.DateStart.Format("2006-01-02"), filter.DateEnd.Format("2006-01-02"))
	} else {
		hargaQuery = hargaQuery.Where("DATE(date_scanned) = ?", time.Now().Format("2006-01-02"))
	}
	hargaQuery.Select("COALESCE(SUM(price),0)").Scan(&result.TotalHargaAsalTerscan)

	return result, nil
}
