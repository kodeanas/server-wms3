package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type BagRepository interface {
	Create(bag *models.Bag) error
	FindByID(id string) (*models.Bag, error)
	FindAll() ([]models.Bag, error)
	FindByType(bagType string) ([]models.Bag, error)
	FindByTypePaginated(bagType string, limit, offset int, search string) ([]models.Bag, int64, error)
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

func (r *bagRepository) FindByType(bagType string) ([]models.Bag, error) {
	var bags []models.Bag
	err := r.db.Where("type = ?", bagType).Find(&bags).Error
	return bags, err
}

func (r *bagRepository) FindByTypePaginated(bagType string, limit, offset int, search string) ([]models.Bag, int64, error) {
	var (
		bags  []models.Bag
		total int64
	)
	query := r.db.Model(&models.Bag{}).Where("type = ?", bagType)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("code ILIKE ?", like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&bags).Error; err != nil {
		return nil, 0, err
	}
	return bags, total, nil
}
