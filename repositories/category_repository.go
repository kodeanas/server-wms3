package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// categoryRepository is GORM implementation.
type categoryRepository struct {
	db *gorm.DB
}

// CategoryRepository defines interface for category CRUD.
type CategoryRepository interface {
	Create(category *models.Category) error
	GetBySlug(slug string) (*models.Category, error)
	GetSlugLike(slug string) ([]models.Category, error)
	List() ([]models.Category, error)
	GetByID(id string) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id string) error
}

func (r *categoryRepository) Delete(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}
func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}
func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	var cat models.Category
	if err := r.db.Where("id = ?", id).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
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
	if err := r.db.Unscoped().Where("slug = ?", slug).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) GetSlugLike(slug string) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Unscoped().Where("slug LIKE ?", slug+"%").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
