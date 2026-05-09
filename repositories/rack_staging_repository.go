package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type RackStagingRepository struct {
	DB *gorm.DB
}

func NewRackStagingRepository(db *gorm.DB) *RackStagingRepository {
	return &RackStagingRepository{DB: db}
}

func (r *RackStagingRepository) Create(rack *models.RackStaging) error {
	return r.DB.Create(rack).Error
}

func (r *RackStagingRepository) FindAll() ([]models.RackStaging, error) {
	var racks []models.RackStaging
	err := r.DB.Where("deleted_at IS NULL").Find(&racks).Error
	return racks, err
}

func (r *RackStagingRepository) FindByID(id string) (*models.RackStaging, error) {
	var rack models.RackStaging
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&rack).Error
	return &rack, err
}

func (r *RackStagingRepository) Update(rack *models.RackStaging) error {
	return r.DB.Save(rack).Error
}

func (r *RackStagingRepository) SoftDelete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&models.RackStaging{}).Error
}

// Get count of rack stagings for a rack display
func (r *RackStagingRepository) CountByRackDisplayID(rackDisplayID string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.RackStaging{}).Where("rack_display_id = ? AND deleted_at IS NULL", rackDisplayID).Count(&count).Error
	return count, err
}

// Find all rack stagings
func (r *RackStagingRepository) FindAllRackStaging() ([]models.RackStaging, error) {
	var racks []models.RackStaging
	err := r.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&racks).Error
	return racks, err
}

// Find all rack stagings with pagination & search
func (r *RackStagingRepository) FindAllRackStagingPaginated(limit, offset int, search string) ([]models.RackStaging, int64, error) {
	var (
		racks []models.RackStaging
		total int64
	)
	query := r.DB.Model(&models.RackStaging{}).Where("deleted_at IS NULL")
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

// Set is_moved = true untuk rack staging tertentu
func (r *RackStagingRepository) SetIsMoved(rackStagingID string) error {
	return r.DB.Model(&models.RackStaging{}).
		Where("id = ? AND deleted_at IS NULL", rackStagingID).
		Update("is_moved", true).Error
}
