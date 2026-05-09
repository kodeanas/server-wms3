package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// stickerRepository is GORM implementation.
type stickerRepository struct {
	db *gorm.DB
}

// StickerRepository defines interface for sticker CRUD.
type StickerRepository interface {
	Create(sticker *models.Sticker) error
	GetBySlug(slug string) (*models.Sticker, error)
	GetSlugLike(slug string) ([]models.Sticker, error)
	GetByID(id string) (*models.Sticker, error)
	List() ([]models.Sticker, error)
	ListPaginated(limit, offset int, search string) ([]models.Sticker, int64, error)
	Update(sticker *models.Sticker) error
	Delete(id string) error
}

// NewStickerRepository constructor.
func NewStickerRepository(db *gorm.DB) StickerRepository {
	return &stickerRepository{db: db}
}

func (r *stickerRepository) Create(sticker *models.Sticker) error {
	return r.db.Create(sticker).Error
}

func (r *stickerRepository) GetBySlug(slug string) (*models.Sticker, error) {
	var sticker models.Sticker
	if err := r.db.Where("slug = ? AND deleted_at IS NULL", slug).First(&sticker).Error; err != nil {
		return nil, err
	}
	return &sticker, nil
}

func (r *stickerRepository) GetSlugLike(slug string) ([]models.Sticker, error) {
	var stickers []models.Sticker
	if err := r.db.Where("slug LIKE ? AND deleted_at IS NULL", slug+"%").Find(&stickers).Error; err != nil {
		return nil, err
	}
	return stickers, nil
}

func (r *stickerRepository) GetByID(id string) (*models.Sticker, error) {
	var sticker models.Sticker
	if err := r.db.Where("id = ?", id).First(&sticker).Error; err != nil {
		return nil, err
	}
	return &sticker, nil
}

func (r *stickerRepository) List() ([]models.Sticker, error) {
	var stickers []models.Sticker
	if err := r.db.Find(&stickers).Error; err != nil {
		return nil, err
	}
	return stickers, nil
}

func (r *stickerRepository) ListPaginated(limit, offset int, search string) ([]models.Sticker, int64, error) {
	var (
		stickers []models.Sticker
		total    int64
	)
	query := r.db.Model(&models.Sticker{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ? OR code_hex ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&stickers).Error; err != nil {
		return nil, 0, err
	}
	return stickers, total, nil
}

func (r *stickerRepository) Update(sticker *models.Sticker) error {
	return r.db.Save(sticker).Error
}

func (r *stickerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Sticker{}).Error
}
