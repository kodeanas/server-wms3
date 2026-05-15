package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// BuyerRepository defines interface for buyer CRUD.
type BuyerRepository interface {
	Create(buyer *models.Buyer) error
	GetByID(id string) (*models.Buyer, error)
	List() ([]models.Buyer, error)
	ListPaginated(limit, offset int, search string) ([]models.Buyer, int64, error)
	ListByClassID(classID string, limit, offset int, search string) ([]models.Buyer, int64, error)
	Update(buyer *models.Buyer) error
	Delete(id string) error
}

// buyerRepository is GORM implementation.
type buyerRepository struct {
	db *gorm.DB
}

// NewBuyerRepository constructor.
func NewBuyerRepository(db *gorm.DB) BuyerRepository {
	return &buyerRepository{db: db}
}

func (r *buyerRepository) Create(buyer *models.Buyer) error {
	return r.db.Create(buyer).Error
}

func (r *buyerRepository) GetByID(id string) (*models.Buyer, error) {
	var buyer models.Buyer
	if err := r.db.Where("id = ?", id).First(&buyer).Error; err != nil {
		return nil, err
	}
	return &buyer, nil
}

func (r *buyerRepository) List() ([]models.Buyer, error) {
	var buyers []models.Buyer
	if err := r.db.Find(&buyers).Error; err != nil {
		return nil, err
	}
	return buyers, nil
}

func (r *buyerRepository) ListPaginated(limit, offset int, search string) ([]models.Buyer, int64, error) {
	var (
		buyers []models.Buyer
		total  int64
	)
	query := r.db.Model(&models.Buyer{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&buyers).Error; err != nil {
		return nil, 0, err
	}
	return buyers, total, nil
}

func (r *buyerRepository) Update(buyer *models.Buyer) error {
	return r.db.Save(buyer).Error
}

func (r *buyerRepository) ListByClassID(classID string, limit, offset int, search string) ([]models.Buyer, int64, error) {
	var (
		buyers []models.Buyer
		total  int64
	)
	query := r.db.Model(&models.Buyer{}).Where("class_id = ?", classID)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&buyers).Error; err != nil {
		return nil, 0, err
	}
	return buyers, total, nil
}

func (r *buyerRepository) Delete(id string) error {
	return r.db.Delete(&models.Buyer{}, "id = ?", id).Error
}
