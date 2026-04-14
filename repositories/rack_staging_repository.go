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
