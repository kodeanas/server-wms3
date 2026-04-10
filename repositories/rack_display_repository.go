package repositories

import (
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
