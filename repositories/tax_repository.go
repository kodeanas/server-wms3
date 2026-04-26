package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type TaxRepository struct {
	DB *gorm.DB
}

func NewTaxRepository(db *gorm.DB) *TaxRepository {
	return &TaxRepository{DB: db}
}

func (r *TaxRepository) Create(tax *models.Tax) error {
	return r.DB.Create(tax).Error
}

func (r *TaxRepository) Update(tax *models.Tax) error {
	return r.DB.Save(tax).Error
}

func (r *TaxRepository) Delete(id string) error {
	return r.DB.Delete(&models.Tax{}, "id = ?", id).Error
}

func (r *TaxRepository) FindByID(id string) (*models.Tax, error) {
	var tax models.Tax
	err := r.DB.First(&tax, "id = ?", id).Error
	return &tax, err
}

func (r *TaxRepository) FindAll() ([]models.Tax, error) {
	var taxes []models.Tax
	err := r.DB.Find(&taxes).Error
	return taxes, err
}

func (r *TaxRepository) SetAllInactiveExcept(id string) error {
	return r.DB.Model(&models.Tax{}).Where("id != ?", id).Update("is_active", false).Error
}

func (r *TaxRepository) FindActive() (*models.Tax, error) {
	var tax models.Tax
	err := r.DB.Where("is_active = ?", true).First(&tax).Error
	return &tax, err
}
