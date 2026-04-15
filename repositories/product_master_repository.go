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
	FindDetailByID(id string) (*dto.ProductMasterDetailResponse, error)
	FindByID(id string) (*models.ProductMaster, error)
	FindStickerByPrice(price float64) (*models.Sticker, error)
	GetCategoryDiscount(categoryID string) (*int, error)
	Update(master *models.ProductMaster) error
	FindByDocumentAndDateRange(documentCode string, from, to time.Time) ([]models.ProductMaster, error)
	Create(master *models.ProductMaster) error
	FindByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	UpdateRackStagingID(id string, rackStagingID string) error
	FindAllByRackStagingID(rackStagingID string) ([]models.ProductMaster, error)
	MoveAllToDisplay(rackStagingID, rackDisplayID string) error
	UpdateBagID(productID string, bagID string) error
	FindByBagID(bagID string) ([]models.ProductMaster, error)
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

func (r *productMasterRepository) Update(master *models.ProductMaster) error {
	return r.db.Save(master).Error
}

func (r *productMasterRepository) FindByID(id string) (*models.ProductMaster, error) {
	var master models.ProductMaster
	err := r.db.Where("id = ?", id).First(&master).Error
	if err != nil {
		return nil, err
	}
	return &master, nil
}

func (r *productMasterRepository) FindStickerByPrice(price float64) (*models.Sticker, error) {
	var sticker models.Sticker
	err := r.db.Where("? BETWEEN min_price AND max_price AND deleted_at IS NULL", price).First(&sticker).Error
	if err != nil {
		return nil, err
	}
	return &sticker, nil
}

func (r *productMasterRepository) GetCategoryDiscount(categoryID string) (*int, error) {
	var result struct {
		Discount *int `gorm:"column:discount"`
	}
	err := r.db.Table("categories").Select("discount").Where("id = ? AND deleted_at IS NULL", categoryID).Take(&result).Error
	if err != nil {
		return nil, err
	}
	return result.Discount, nil
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

type productMasterDetailRow struct {
	ID               string    `gorm:"column:id"`
	DocumentID       string    `gorm:"column:document_id"`
	Barcode          string    `gorm:"column:barcode"`
	BarcodeWarehouse string    `gorm:"column:barcode_warehouse"`
	Name             string    `gorm:"column:name"`
	NameWarehouse    string    `gorm:"column:name_warehouse"`
	Item             int       `gorm:"column:item"`
	ItemWarehouse    int       `gorm:"column:item_warehouse"`
	Price            float64   `gorm:"column:price"`
	PriceWarehouse   float64   `gorm:"column:price_warehouse"`
	CategoryID       *string   `gorm:"column:category_id"`
	StickerID        *string   `gorm:"column:sticker_id"`
	TypeOut          *string   `gorm:"column:type_out"`
	Location         string    `gorm:"column:location"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	DocumentCode     *string   `gorm:"column:document_code"`
	DocumentFileName *string   `gorm:"column:document_file_name"`
	CategoryName     *string   `gorm:"column:category_name"`
	CategoryDiscount *int      `gorm:"column:category_discount"`
	StickerName      *string   `gorm:"column:sticker_name"`
	StickerCodeHex   *string   `gorm:"column:sticker_code_hex"`
	StickerType      *string   `gorm:"column:sticker_type"`
}

func (r *productMasterRepository) FindDetailByID(id string) (*dto.ProductMasterDetailResponse, error) {
	var row productMasterDetailRow
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
			pm.sticker_id,
			pm.type_out,
			pm.location,
			pm.created_at,
			pd.code AS document_code,
			pd.file_name AS document_file_name,
			c.name AS category_name,
			c.discount AS category_discount,
			st.name AS sticker_name,
			st.code_hex AS sticker_code_hex,
			st.type AS sticker_type
		`).
		Joins("LEFT JOIN product_documents pd ON pd.id = pm.document_id::uuid").
		Joins("LEFT JOIN categories c ON c.id = pm.category_id::uuid").
		Joins("LEFT JOIN stickers st ON st.id = pm.sticker_id::uuid").
		Where("pm.id = ?", id).
		Take(&row).Error
	if err != nil {
		return nil, err
	}

	detail := &dto.ProductMasterDetailResponse{
		ID:               row.ID,
		DocumentID:       row.DocumentID,
		DocumentCode:     derefString(row.DocumentCode),
		DocumentName:     derefString(row.DocumentFileName),
		Barcode:          row.Barcode,
		BarcodeWarehouse: row.BarcodeWarehouse,
		Name:             row.Name,
		NameWarehouse:    row.NameWarehouse,
		Item:             row.Item,
		ItemWarehouse:    row.ItemWarehouse,
		Price:            row.Price,
		PriceWarehouse:   row.PriceWarehouse,
		CategoryID:       row.CategoryID,
		StickerID:        row.StickerID,
		TypeOut:          displayTypeOut(row.TypeOut),
		Location:         row.Location,
		CreatedAt:        row.CreatedAt,
		Additional: dto.ProductMasterDetailAdditional{
			Document: nil,
			Sticker:  nil,
			Category: nil,
		},
	}

	if row.DocumentCode != nil || row.DocumentFileName != nil {
		detail.Additional.Document = &dto.ProductMasterDocumentAdditional{
			Code:     derefString(row.DocumentCode),
			NameFile: derefString(row.DocumentFileName),
		}
	}

	// Logika: jika sticker ada, prioritaskan sticker dan set category null
	// jika sticker null tapi category ada, gunakan category dan set sticker null
	if row.StickerName != nil {
		detail.Additional.Sticker = &dto.ProductMasterStickerAdditional{
			Name:    *row.StickerName,
			CodeHex: derefString(row.StickerCodeHex),
			Type:    derefString(row.StickerType),
		}
		detail.Additional.Category = nil // pastikan category null jika sticker ada
	} else if row.CategoryName != nil {
		detail.Additional.Category = &dto.ProductMasterCategoryAdditional{
			Name:    *row.CategoryName,
			Diskon:  derefInt(row.CategoryDiscount),
			Type:    "",
			CodeHex: "",
		}
		detail.Additional.Sticker = nil // pastikan sticker null jika category ada
	}

	return detail, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func displayTypeOut(value *string) string {
	if value == nil || *value == "" {
		return "belum keluar"
	}
	return *value
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
			"location":        "display",
		}).Error
}

// Update bag_id untuk product master
func (r *productMasterRepository) UpdateBagID(productID string, bagID string) error {
	return r.db.Model(&models.ProductMaster{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("bag_id", bagID).Error
}

// List all product master in a bag (rackStagingSticker)
func (r *productMasterRepository) FindByBagID(bagID string) ([]models.ProductMaster, error) {
	var masters []models.ProductMaster
	err := r.db.Where("bag_id = ? AND deleted_at IS NULL", bagID).Order("created_at DESC").Find(&masters).Error
	return masters, err
}
