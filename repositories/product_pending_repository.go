package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type ProductPendingRepository interface {
	FindByDocumentID(documentID string) ([]models.ProductPending, error)
	FindByDocumentIDAndBarcode(documentID, barcode string) (*models.ProductPending, error)
	Update(product *models.ProductPending) error
}

type productPendingRepository struct {
	db *gorm.DB
}

func NewProductPendingRepository(db *gorm.DB) ProductPendingRepository {
	return &productPendingRepository{db: db}
}

func (r *productPendingRepository) FindByDocumentID(documentID string) ([]models.ProductPending, error) {
	var products []models.ProductPending
	err := r.db.Where("document_id = ? AND deleted_at IS NULL", documentID).Find(&products).Error
	return products, err
}

func (r *productPendingRepository) FindByDocumentIDAndBarcode(documentID, barcode string) (*models.ProductPending, error) {
	var product models.ProductPending
	err := r.db.Where("document_id = ? AND barcode = ? AND deleted_at IS NULL", documentID, barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productPendingRepository) Update(product *models.ProductPending) error {
	return r.db.Save(product).Error
}
