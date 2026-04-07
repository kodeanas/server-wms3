package repositories

import (
	"time"
	dto "wms/dto/response"
	"wms/models"

	"gorm.io/gorm"
)

type ProductMasterRepository interface {
	FindByLocation(location string) ([]models.ProductMaster, error)
	FindStagingReguler() ([]dto.ProductMasterRegulerResponse, error)
	FindStagingSticker() ([]dto.ProductMasterStickerResponse, error)
	FindByDocumentAndDateRange(documentCode string, from, to time.Time) ([]models.ProductMaster, error)
}

type productMasterRepository struct {
	db *gorm.DB
}

func NewProductMasterRepository(db *gorm.DB) ProductMasterRepository {
	return &productMasterRepository{db: db}
}

func (r *productMasterRepository) FindByLocation(location string) ([]models.ProductMaster, error) {
	var masters []models.ProductMaster
	err := r.db.Where("location = ?", location).Order("created_at DESC").Find(&masters).Error
	return masters, err
}

func (r *productMasterRepository) FindStagingReguler() ([]dto.ProductMasterRegulerResponse, error) {
	var masters []dto.ProductMasterRegulerResponse
	err := r.db.Table("product_masters pm").
		Select(`
			pm.id,
			pm.document_id,
			pm.barcode,
			pm.barcode_warehouse,
			pm.name,
			pm.name_warehouse,
			pm.item,
			pm.item_warehouse,
			pm.price,
			pm.price_warehouse,
			pm.category_id,
			pm.product_pending_id,
			pm.is_sku,
			pm.location,
			pm.type_out,
			pm.created_at,
			pm.updated_at,
			categories.name AS category_name
		`).
		Joins("LEFT JOIN categories ON categories.id = pm.category_id::uuid").
		Where("pm.location = ?", "staging_reguler").
		Order("pm.created_at DESC").
		Scan(&masters).Error
	return masters, err
}

func (r *productMasterRepository) FindStagingSticker() ([]dto.ProductMasterStickerResponse, error) {
	var masters []dto.ProductMasterStickerResponse
	err := r.db.Table("product_masters pm").
		Select(`
			pm.id,
			pm.document_id,
			pm.barcode,
			pm.barcode_warehouse,
			pm.name,
			pm.name_warehouse,
			pm.item,
			pm.item_warehouse,
			pm.price,
			pm.price_warehouse,
			pm.sticker_id,
			pm.product_pending_id,
			pm.is_sku,
			pm.location,
			pm.type_out,
			pm.created_at,
			pm.updated_at,
			stickers.name AS sticker_name,
			stickers.code_hex AS sticker_code_hex
		`).
		Joins("LEFT JOIN stickers ON stickers.id = pm.sticker_id::uuid").
		Where("pm.location = ?", "staging_sticker").
		Order("pm.created_at DESC").
		Scan(&masters).Error
	return masters, err
}

func (r *productMasterRepository) FindByDocumentAndDateRange(documentCode string, from, to time.Time) ([]models.ProductMaster, error) {
	var masters []models.ProductMaster
	err := r.db.Joins("JOIN product_documents ON product_documents.id = product_masters.document_id::uuid").
		Where("product_documents.code = ? AND product_masters.created_at BETWEEN ? AND ?", documentCode, from, to).
		Order("product_masters.created_at DESC").
		Find(&masters).Error
	return masters, err
}
