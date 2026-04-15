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
	Create(master *models.ProductMaster) error
	FindByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	UpdateRackStagingID(id string, rackStagingID string) error
	FindAllByRackStagingID(rackStagingID string) ([]models.ProductMaster, error)
	MoveAllToDisplay(rackStagingID, rackDisplayID string) error
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

func (r *productMasterRepository) Create(master *models.ProductMaster) error {
	return r.db.Create(master).Error
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

// Find product master by barcode_warehouse
func (r *productMasterRepository) FindByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	var master models.ProductMaster
	err := r.db.Where("barcode_warehouse = ? AND deleted_at IS NULL", barcode).First(&master).Error
	if err != nil {
		return nil, err
	}
	return &master, nil
}

// Update rack_staging_id for product master
func (r *productMasterRepository) UpdateRackStagingID(id string, rackStagingID string) error {
	return r.db.Model(&models.ProductMaster{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("rack_staging_id", rackStagingID).Error
}

// Find all product masters by rack staging id
func (r *productMasterRepository) FindAllByRackStagingID(rackStagingID string) ([]models.ProductMaster, error) {
	var masters []models.ProductMaster
	err := r.db.Where("rack_staging_id = ? AND deleted_at IS NULL", rackStagingID).Order("created_at DESC").Find(&masters).Error
	return masters, err
}

// Update massal product master pada rack staging: set rack_display_id dan location
func (r *productMasterRepository) MoveAllToDisplay(rackStagingID, rackDisplayID string) error {
	return r.db.Model(&models.ProductMaster{}).
		Where("rack_staging_id = ? AND deleted_at IS NULL", rackStagingID).
		Updates(map[string]interface{}{
			"rack_display_id": rackDisplayID,
			"location": "display",
		}).Error
}
