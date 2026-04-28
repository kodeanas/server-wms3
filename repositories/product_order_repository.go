package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type ProductOrderRepository interface {
	Create(productOrder *models.ProductOrder) error
	GetByID(id string) (*models.ProductOrder, error)
	ListByOrderID(orderID string) ([]models.ProductOrder, error)
	Update(productOrder *models.ProductOrder) error
	Delete(id string) error
}

type productOrderRepository struct {
	db *gorm.DB
}

func NewProductOrderRepository(db *gorm.DB) ProductOrderRepository {
	return &productOrderRepository{db: db}
}

func (r *productOrderRepository) Create(productOrder *models.ProductOrder) error {
	return r.db.Create(productOrder).Error
}

func (r *productOrderRepository) GetByID(id string) (*models.ProductOrder, error) {
	var po models.ProductOrder
	if err := r.db.Where("id = ?", id).First(&po).Error; err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *productOrderRepository) ListByOrderID(orderID string) ([]models.ProductOrder, error) {
	var pos []models.ProductOrder
	if err := r.db.Where("order_id = ?", orderID).Find(&pos).Error; err != nil {
		return nil, err
	}
	return pos, nil
}

func (r *productOrderRepository) Update(productOrder *models.ProductOrder) error {
	return r.db.Save(productOrder).Error
}

func (r *productOrderRepository) Delete(id string) error {
	return r.db.Delete(&models.ProductOrder{}, "id = ?", id).Error
}
