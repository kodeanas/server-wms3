package repositories

import (
	"time"
	"wms/models"

	"gorm.io/gorm"
)

type ProductDocumentRepository interface {
	FindAll() ([]models.ProductDocument, error)
	FindAllPaginated(limit, offset int, search string) ([]models.ProductDocument, int64, error)
	Create(doc *models.ProductDocument) error
	FindByType(docType string) ([]models.ProductDocument, error)
	FindByTypePaginated(docType string, limit, offset int, search string) ([]models.ProductDocument, int64, error)
	FindBulkDetailByID(id string) (models.ProductDocument, error)
	FindBastByID(id string) (models.ProductDocument, error)
	FindBastRelationsByID(id string) (models.ProductDocument, error)
	FindBastProductPendingByDiscrepancy(id string) ([]models.ProductPending, error)
	FindBastProductPendingByNonDiscrepancy(id string) ([]models.ProductPending, error)
	FindBastScannedSummary(id string) (totalItemScanned int64, totalPriceScanned float64, err error)
	FindBastPendingSummaryByStatuses(id string, statuses []string) (map[string]map[string]float64, error)
	FindSkuDetailByID(id string) (models.ProductDocument, error)
	UpdateDateStopByID(id string, dateStop *time.Time) error
	FindBastDetailByID(id string) (models.ProductDocument, error)
	UpdateStatusByID(id string, status string) error
	GetBastSummaryAll() (map[string]interface{}, error)
	GetBulkSummaryAll() (map[string]interface{}, error)
	GetSkuSummaryAll() (map[string]interface{}, error)
}

// UpdateStatusByID mengubah status dokumen
func (r *productDocumentRepository) UpdateStatusByID(id string, status string) error {
	return r.db.Model(&models.ProductDocument{}).Where("id = ?", id).Update("status", status).Error
}

type productDocumentRepository struct {
	db *gorm.DB
}

func NewProductDocumentRepository(db *gorm.DB) ProductDocumentRepository {
	return &productDocumentRepository{db: db}
}

func (r *productDocumentRepository) Create(doc *models.ProductDocument) error {
	return r.db.Create(doc).Error
}

func (r *productDocumentRepository) FindByType(docType string) ([]models.ProductDocument, error) {
	var documents []models.ProductDocument
	// Mengambil data berdasarkan type (misal: 'bulk')
	err := r.db.Where("type = ?", docType).Find(&documents).Error
	return documents, err
}

func (r *productDocumentRepository) FindByTypePaginated(docType string, limit, offset int, search string) ([]models.ProductDocument, int64, error) {
	var (
		documents []models.ProductDocument
		total     int64
	)
	query := r.db.Model(&models.ProductDocument{}).Where("type = ?", docType)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("code ILIKE ? OR file_name ILIKE ? OR supplier ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&documents).Error; err != nil {
		return nil, 0, err
	}
	return documents, total, nil
}

func (r *productDocumentRepository) FindAll() ([]models.ProductDocument, error) {
	var docs []models.ProductDocument
	err := r.db.Order("created_at DESC").Find(&docs).Error
	return docs, err
}

func (r *productDocumentRepository) FindAllPaginated(limit, offset int, search string) ([]models.ProductDocument, int64, error) {
	var (
		documents []models.ProductDocument
		total     int64
	)
	query := r.db.Model(&models.ProductDocument{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("code ILIKE ? OR file_name ILIKE ? OR supplier ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&documents).Error; err != nil {
		return nil, 0, err
	}
	return documents, total, nil
}

func (r *productDocumentRepository) FindSkuDetailByID(id string) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := r.db.Preload("ProductPendings").
		Where("id = ? AND type = ?", id, "sku").
		First(&doc).Error
	return doc, err
}

func (r *productDocumentRepository) FindBastDetailByID(id string) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := r.db.Preload("ProductPendings").
		Where("id = ? AND type = ?", id, "bast").
		First(&doc).Error
	return doc, err
}
func (r *productDocumentRepository) FindBulkDetailByID(id string) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := r.db.Preload("ProductPendings").
		Where("id = ? AND type = ?", id, "bulk").
		First(&doc).Error
	return doc, err
}

func (r *productDocumentRepository) FindBastByID(id string) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := r.db.Where("id = ? AND type = ?", id, "bast").First(&doc).Error
	return doc, err
}

func (r *productDocumentRepository) FindBastRelationsByID(id string) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := r.db.Preload("ProductMasters").
		Preload("ProductPendings").
		Where("id = ? AND type = ?", id, "bast").
		First(&doc).Error
	return doc, err
}

func (r *productDocumentRepository) FindBastProductPendingByDiscrepancy(id string) ([]models.ProductPending, error) {
	var pendings []models.ProductPending
	err := r.db.Model(&models.ProductPending{}).
		Where("document_id = ? AND status = ?", id, "discrepancy").
		Order("created_at DESC").
		Find(&pendings).Error
	return pendings, err
}

func (r *productDocumentRepository) FindBastProductPendingByNonDiscrepancy(id string) ([]models.ProductPending, error) {
	var pendings []models.ProductPending
	err := r.db.Model(&models.ProductPending{}).
		Where("document_id = ? AND status <> ?", id, "discrepancy").
		Order("date_scanned DESC NULLS LAST").
		Order("created_at DESC").
		Find(&pendings).Error
	return pendings, err
}

func (r *productDocumentRepository) FindBastScannedSummary(id string) (totalItemScanned int64, totalPriceScanned float64, err error) {
	err = r.db.Model(&models.ProductPending{}).
		Where("document_id = ? AND date_scanned IS NOT NULL", id).
		Count(&totalItemScanned).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.ProductPending{}).
		Where("document_id = ? AND date_scanned IS NOT NULL", id).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalPriceScanned).Error
	if err != nil {
		return 0, 0, err
	}

	return totalItemScanned, totalPriceScanned, nil
}

func (r *productDocumentRepository) FindBastPendingSummaryByStatuses(id string, statuses []string) (map[string]map[string]float64, error) {
	type statusSummaryRow struct {
		Status     string
		TotalItem  int64
		TotalPrice float64
	}

	var rows []statusSummaryRow
	err := r.db.Model(&models.ProductPending{}).
		Select("status, COUNT(*) AS total_item, COALESCE(SUM(price), 0) AS total_price").
		Where("document_id = ? AND status IN ?", id, statuses).
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]float64)
	for _, status := range statuses {
		result[status] = map[string]float64{
			"total_item":  0,
			"total_price": 0,
		}
	}

	for _, row := range rows {
		result[row.Status] = map[string]float64{
			"total_item":  float64(row.TotalItem),
			"total_price": row.TotalPrice,
		}
	}

	return result, nil
}

// UpdateDateStopByID mengisi field date_stop pada dokumen
func (r *productDocumentRepository) UpdateDateStopByID(id string, dateStop *time.Time) error {
	return r.db.Model(&models.ProductDocument{}).Where("id = ?", id).Update("date_stop", dateStop).Error
}

// GetBastSummaryAll returns summary all for BAST with total inbound, scanned, and product status breakdown
func (r *productDocumentRepository) GetBastSummaryAll() (map[string]interface{}, error) {
	// Pakai timezone Asia/Jakarta agar "hari ini" konsisten dengan waktu lokal user.
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Count total documents inbound today (type = 'bast')
	var totalDocInbound int64
	if err := r.db.Model(&models.ProductDocument{}).
		Where("type = ? AND created_at >= ? AND created_at < ?", "bast", startOfDay, endOfDay).
		Count(&totalDocInbound).Error; err != nil {
		return nil, err
	}

	// Count total documents scanned today (type = 'bast' and status != 'progress')
	var totalDocScanned int64
	if err := r.db.Model(&models.ProductDocument{}).
		Where("type = ? AND status != ? AND date_stop >= ? AND date_stop < ?", "bast", "progress", startOfDay, endOfDay).
		Count(&totalDocScanned).Error; err != nil {
		return nil, err
	}

	// Get product status breakdown from product_pending (not deleted and has document with type bast)
	var statusCounts struct {
		Good     int64 `gorm:"column:good_count"`
		Damaged  int64 `gorm:"column:damaged_count"`
		Abnormal int64 `gorm:"column:abnormal_count"`
		Non      int64 `gorm:"column:non_count"`
	}

	query := r.db.Table("product_pendings pp").
		Select(`
			COUNT(CASE WHEN pp.status = 'good' THEN 1 END) as good_count,
			COUNT(CASE WHEN pp.status = 'damaged' THEN 1 END) as damaged_count,
			COUNT(CASE WHEN pp.status = 'abnormal' THEN 1 END) as abnormal_count,
			COUNT(CASE WHEN pp.status = 'non' THEN 1 END) as non_count
		`).
		Joins("INNER JOIN product_documents pd ON pp.document_id::text = pd.id::text").
		Where("pp.deleted_at IS NULL").
		Where("pd.deleted_at IS NULL").
		Where("pd.type = ?", "bast")

	if err := query.Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_document_inbound": totalDocInbound,
		"total_document_scanned": totalDocScanned,
		"total_product_good":     statusCounts.Good,
		"total_product_damaged":  statusCounts.Damaged,
		"total_product_abnormal": statusCounts.Abnormal,
		"total_product_non":      statusCounts.Non,
	}, nil
}

// GetBulkSummaryAll returns summary all for BULK with total document upload, product masuk, and harga masuk
func (r *productDocumentRepository) GetBulkSummaryAll() (map[string]interface{}, error) {
	// Total dokumen ter-upload (type = 'bulk')
	var totalDocUpload int64
	if err := r.db.Model(&models.ProductDocument{}).
		Where("type = ?", "bulk").
		Count(&totalDocUpload).Error; err != nil {
		return nil, err
	}

	// Total product masuk: dari product_pendings yang document.type = 'bulk'
	var totalProductMasuk int64
	if err := r.db.Table("product_pendings pp").
		Joins("INNER JOIN product_documents pd ON pp.document_id::text = pd.id::text").
		Where("pp.deleted_at IS NULL").
		Where("pd.deleted_at IS NULL").
		Where("pd.type = ?", "bulk").
		Count(&totalProductMasuk).Error; err != nil {
		return nil, err
	}

	// Total harga masuk: SUM(price_warehouse) dari product_masters yang document.type = 'bulk'
	var totalHargaMasuk float64
	if err := r.db.Table("product_masters pm").
		Select("COALESCE(SUM(pm.price_warehouse), 0)").
		Joins("INNER JOIN product_documents pd ON pm.document_id::text = pd.id::text").
		Where("pm.deleted_at IS NULL").
		Where("pd.deleted_at IS NULL").
		Where("pd.type = ?", "bulk").
		Scan(&totalHargaMasuk).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_document_upload": totalDocUpload,
		"total_product_masuk":   totalProductMasuk,
		"total_harga_masuk":     totalHargaMasuk,
	}, nil
}

// GetSkuSummaryAll returns summary all for SKU with total document upload, product masuk, and harga masuk
func (r *productDocumentRepository) GetSkuSummaryAll() (map[string]interface{}, error) {
	// Total dokumen ter-upload (type = 'sku')
	var totalDocUpload int64
	if err := r.db.Model(&models.ProductDocument{}).
		Where("type = ?", "sku").
		Count(&totalDocUpload).Error; err != nil {
		return nil, err
	}

	// Total product masuk: dari product_pendings yang document.type = 'sku'
	var totalProductMasuk int64
	if err := r.db.Table("product_pendings pp").
		Joins("INNER JOIN product_documents pd ON pp.document_id::text = pd.id::text").
		Where("pp.deleted_at IS NULL").
		Where("pd.deleted_at IS NULL").
		Where("pd.type = ?", "sku").
		Count(&totalProductMasuk).Error; err != nil {
		return nil, err
	}

	// Total harga masuk: SUM(price) dari product_pendings yang document.type = 'sku'
	var totalHargaMasuk float64
	if err := r.db.Table("product_pendings pp").
		Select("COALESCE(SUM(pp.price), 0)").
		Joins("INNER JOIN product_documents pd ON pp.document_id::text = pd.id::text").
		Where("pp.deleted_at IS NULL").
		Where("pd.deleted_at IS NULL").
		Where("pd.type = ?", "sku").
		Scan(&totalHargaMasuk).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_document_upload": totalDocUpload,
		"total_product_masuk":   totalProductMasuk,
		"total_harga_masuk":     totalHargaMasuk,
	}, nil
}
