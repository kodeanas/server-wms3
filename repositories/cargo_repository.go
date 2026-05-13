package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type CargoRepository interface {
	Create(cargo *models.Cargo) error
	FindByID(id string) (*models.Cargo, error)
	FindAll() ([]models.Cargo, error)
	FindAllPaginated(limit, offset int, search string) ([]models.Cargo, int64, error)
	SetStatus(cargoID, status string) error
	SetIsSale(cargoID string, isSale bool) error
	SetIsMoved(cargoID string) error
}

type cargoRepository struct {
	db *gorm.DB
}

func NewCargoRepository(db *gorm.DB) CargoRepository {
	return &cargoRepository{db: db}
}

func (r *cargoRepository) Create(cargo *models.Cargo) error {
	return r.db.Create(cargo).Error
}

func (r *cargoRepository) FindByID(id string) (*models.Cargo, error) {
	var cargo models.Cargo
	err := r.db.Where("id = ?", id).First(&cargo).Error

	return &cargo, err
}

func (r *cargoRepository) FindAll() ([]models.Cargo, error) {
	var cargo []models.Cargo
	err := r.db.Find(&cargo).Error

	return cargo, err
}

func (r *cargoRepository) FindAllPaginated(limit, offset int, search string) ([]models.Cargo, int64, error) {
	var (
		cargos []models.Cargo
		total  int64
	)
	query := r.db.Model(&models.Cargo{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("code ILIKE ?", like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(limit).Offset(offset).Find(&cargos).Error
	return cargos, total, err
}

func (r *cargoRepository) SetStatus(cargoID, status string) error {
	return r.db.Model(&models.Cargo{}).Where("id = ?", cargoID).Update("status", status).Error
}

func (r *cargoRepository) SetIsSale(cargoID string, isSale bool) error {
	return r.db.Model(&models.Cargo{}).Where("id = ?", cargoID).Update("is_sale", isSale).Error
}

func (r *cargoRepository) SetIsMoved(cargoID string) error {
	return r.db.Model(&models.Bag{}).Where("cargo_id = ?", cargoID).Update("is_moved", true).Error
}
