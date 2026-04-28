package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type DiscountOrderRepository interface {
	Create(discountOrder *models.DiscountOrder) error
	GetByID(id string) (*models.DiscountOrder, error)
	ListByOrderID(orderID string) ([]models.DiscountOrder, error)
	Update(discountOrder *models.DiscountOrder) error
	Delete(id string) error
}

type discountOrderRepository struct {
	db *gorm.DB
}

func NewDiscountOrderRepository(db *gorm.DB) DiscountOrderRepository {
	return &discountOrderRepository{db: db}
}

func (r *discountOrderRepository) Create(discountOrder *models.DiscountOrder) error {
	return r.db.Create(discountOrder).Error
}

func (r *discountOrderRepository) GetByID(id string) (*models.DiscountOrder, error) {
	var do models.DiscountOrder
	if err := r.db.Where("id = ?", id).First(&do).Error; err != nil {
		return nil, err
	}
	return &do, nil
}

func (r *discountOrderRepository) ListByOrderID(orderID string) ([]models.DiscountOrder, error) {
	var dos []models.DiscountOrder
	if err := r.db.Where("order_id = ?", orderID).Find(&dos).Error; err != nil {
		return nil, err
	}
	return dos, nil
}

func (r *discountOrderRepository) Update(discountOrder *models.DiscountOrder) error {
	return r.db.Save(discountOrder).Error
}

func (r *discountOrderRepository) Delete(id string) error {
	return r.db.Delete(&models.DiscountOrder{}, "id = ?", id).Error
}
