package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type ProductRepairRepository interface {
	Create(repair *models.ProductRepair) error
}

type productRepairRepository struct {
	db *gorm.DB
}

func NewProductRepairRepository(db *gorm.DB) ProductRepairRepository {
	return &productRepairRepository{db: db}
}

func (r *productRepairRepository) Create(repair *models.ProductRepair) error {
	return r.db.Create(repair).Error
}
