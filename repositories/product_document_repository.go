package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type ProductDocumentRepository interface {
	FindAll() ([]models.ProductDocument, error)
	// Tambahkan baris di bawah ini:
	FindByType(docType string) ([]models.ProductDocument, error)
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
