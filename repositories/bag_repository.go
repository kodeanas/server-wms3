package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type BagRepository interface {
	Create(bag *models.Bag) error
	FindByID(id string) (*models.Bag, error)
	FindAll() ([]models.Bag, error)
}

type bagRepository struct {
	db *gorm.DB
}

func NewBagRepository(db *gorm.DB) BagRepository {
	return &bagRepository{db: db}
}

func (r *bagRepository) Create(bag *models.Bag) error {
	return r.db.Create(bag).Error
}

func (r *bagRepository) FindByID(id string) (*models.Bag, error) {
	var bag models.Bag
	err := r.db.Where("id = ?", id).First(&bag).Error
	if err != nil {
		return nil, err
	}
	return &bag, nil
}

func (r *bagRepository) FindAll() ([]models.Bag, error) {
	var bags []models.Bag
	err := r.db.Find(&bags).Error
	return bags, err
}
