package repositories

import (
	dto "wms/dto/response"
	"wms/models"

	"gorm.io/gorm"
)

type RackDisplayRepository struct {
	DB *gorm.DB
}

func NewRackDisplayRepository(db *gorm.DB) *RackDisplayRepository {
	return &RackDisplayRepository{DB: db}
}

func (r *RackDisplayRepository) Create(rack *models.RackDisplay) error {
	return r.DB.Create(rack).Error
}

func (r *RackDisplayRepository) FindAll() ([]models.RackDisplay, error) {
	var racks []models.RackDisplay
	err := r.DB.Where("deleted_at IS NULL").Find(&racks).Error
	return racks, err
}

func (r *RackDisplayRepository) FindAllPaginated(limit, offset int, search string) ([]models.RackDisplay, int64, error) {
	var (
		racks []models.RackDisplay
		total int64
	)
	query := r.DB.Model(&models.RackDisplay{}).Where("deleted_at IS NULL")
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR code ILIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&racks).Error; err != nil {
		return nil, 0, err
	}
	return racks, total, nil
}

func (r *RackDisplayRepository) FindByID(id string) (*models.RackDisplay, error) {
	var rack models.RackDisplay
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&rack).Error
	return &rack, err
}

func (r *RackDisplayRepository) Update(rack *models.RackDisplay) error {
	return r.DB.Save(rack).Error
}

func (r *RackDisplayRepository) SoftDelete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&models.RackDisplay{}).Error
}

func (r *RackDisplayRepository) GetSummaryAllDisplay() (*dto.RackDisplaySummaryAllResponse, error) {
	var summary dto.RackDisplaySummaryAllResponse

	err := r.DB.Table("product_masters pm").
		Select(`
			COUNT(DISTINCT pm.rack_display_id) AS total_rack,
			COUNT(pm.id) AS total_item,
			COALESCE(SUM(pm.price), 0) AS total_price
		`).
		Where("pm.deleted_at IS NULL").
		Where("pm.location = ?", "display").
		Where("pm.rack_display_id IS NOT NULL").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	summary.TotalItemAndPrice = float64(summary.TotalItem) + summary.TotalPrice
	return &summary, nil
}

func (r *RackDisplayRepository) GetRackProductSummaryAll() (*dto.RackProductSummaryResponse, error) {
	var summary dto.RackProductSummaryResponse

	err := r.DB.Table("product_masters pm").
		Select(`
			COUNT(pm.id) AS total_item,
			COALESCE(SUM(pm.price), 0) AS total_price,
			COALESCE(SUM(pm.price_warehouse), 0) AS total_price_warehouse
		`).
		Where("pm.deleted_at IS NULL").
		Where("pm.location = ?", "display").
		Where("pm.rack_display_id IS NOT NULL").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	summary.TotalItemAndPrice = float64(summary.TotalItem) + summary.TotalPrice + summary.TotalPriceWarehouse
	return &summary, nil
}
