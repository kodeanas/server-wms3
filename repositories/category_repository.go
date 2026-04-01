package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// CategoryRepository defines interface for category CRUD.
type CategoryRepository interface {
	Create(category *models.Category) error
	GetBySlug(slug string) (*models.Category, error)
	List() ([]models.Category, error)
}

// categoryRepository is GORM implementation.
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository constructor.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var cat models.Category
	if err := r.db.Where("slug = ?", slug).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
