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
	ListPaginated(limit, offset int, search string) ([]models.Category, int64, error)
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

func (r *categoryRepository) ListPaginated(limit, offset int, search string) ([]models.Category, int64, error) {
	var (
		categories []models.Category
		total      int64
	)
	query := r.db.Model(&models.Category{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}
