package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	Update(order *models.Order) error
	Delete(id string) error
	CountDoneByBuyer(buyerID string) (int, error)
	GetLastDoneByBuyer(buyerID string) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id string) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *orderRepository) CountDoneByBuyer(buyerID string) (int, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("buyer_id = ? AND status = ?", buyerID, "done").Count(&count).Error
	return int(count), err
}

func (r *orderRepository) GetLastDoneByBuyer(buyerID string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("buyer_id = ? AND status = ?", buyerID, "done").Order("updated_at DESC").First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
