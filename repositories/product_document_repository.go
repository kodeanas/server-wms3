package repositories

import (
	"time"
	"wms/models"

	"gorm.io/gorm"
)

type ProductDocumentRepository interface {
	FindAll() ([]models.ProductDocument, error)
	// Tambahkan baris di bawah ini:
	FindByType(docType string) ([]models.ProductDocument, error)
	FindBulkDetailByID(id string) (models.ProductDocument, error)
	FindBastByID(id string) (models.ProductDocument, error)
	FindBastRelationsByID(id string) (models.ProductDocument, error)
	FindBastProductPendingByDiscrepancy(id string) ([]models.ProductPending, error)
	FindBastProductPendingByNonDiscrepancy(id string) ([]models.ProductPending, error)
	FindBastScannedSummary(id string) (totalItemScanned int64, totalPriceScanned float64, err error)
	FindBastPendingSummaryByStatuses(id string, statuses []string) (map[string]map[string]float64, error)
	// UpdateDateStopByID mengisi field date_stop pada dokumen
	UpdateDateStopByID(id string, dateStop *time.Time) error
	// UpdateStatusByID mengubah status dokumen
	UpdateStatusByID(id string, status string) error
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

func (r *productDocumentRepository) FindByType(docType string) ([]models.ProductDocument, error) {
	var documents []models.ProductDocument
	// Mengambil data berdasarkan type (misal: 'bulk')
	err := r.db.Where("type = ?", docType).Find(&documents).Error
	return documents, err
}

func (r *productDocumentRepository) FindAll() ([]models.ProductDocument, error) {
	var docs []models.ProductDocument
	err := r.db.Order("created_at DESC").Find(&docs).Error
	return docs, err
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
