package repositories

import (
	"context"
	"wms/config"
	"wms/models"
)

type orderRepository struct{}

// NewOrderRepository creates a new order repository
func NewOrderRepository() OrderRepository {
	return &orderRepository{}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return config.DB.WithContext(ctx).Create(order).Error
}

// GetByID gets order by ID
func (r *orderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := config.DB.WithContext(ctx).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByCode gets order by code
func (r *orderRepository) GetByCode(ctx context.Context, code string) (*models.Order, error) {
	var order models.Order
	err := config.DB.WithContext(ctx).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		First(&order, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID gets orders by user
func (r *orderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	err := config.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetByBuyerID gets orders by buyer
func (r *orderRepository) GetByBuyerID(ctx context.Context, buyerID string, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	err := config.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("buyer_id = ?", buyerID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Where("buyer_id = ?", buyerID).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetAll gets all orders
func (r *orderRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	err := config.DB.WithContext(ctx).Model(&models.Order{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Update updates order
func (r *orderRepository) Update(ctx context.Context, id string, order *models.Order) error {
	return config.DB.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Updates(order).Error
}

// Delete deletes order
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return config.DB.WithContext(ctx).Delete(&models.Order{}, "id = ?", id).Error
}

// GetByStatus gets orders by status
func (r *orderRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	err := config.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ?", status).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Where("status = ?", status).
		Preload("User").
		Preload("Buyer").
		Preload("Items").
		Preload("Cargos").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
