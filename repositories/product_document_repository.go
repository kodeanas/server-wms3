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

// UpdateDateStopByID mengisi field date_stop pada dokumen
func (r *productDocumentRepository) UpdateDateStopByID(id string, dateStop *time.Time) error {
	return r.db.Model(&models.ProductDocument{}).Where("id = ?", id).Update("date_stop", dateStop).Error
}
