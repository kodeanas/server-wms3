package services

import (
	"context"
	"wms/models"
	"wms/repositories"
)

// OrderService defines order business logic
type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	GetOrderByCode(ctx context.Context, code string) (*models.Order, error)
	GetOrdersByUser(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error)
	GetOrdersByBuyer(ctx context.Context, buyerID string, limit, offset int) ([]models.Order, int64, error)
	GetOrdersByStatus(ctx context.Context, status string, limit, offset int) ([]models.Order, int64, error)
	ListOrders(ctx context.Context, limit, offset int) ([]models.Order, int64, error)
	UpdateOrder(ctx context.Context, id string, order *models.Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderStats(ctx context.Context) (map[string]interface{}, error)
}

type orderService struct {
	repo repositories.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

// CreateOrder creates a new order
func (s *orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	return s.repo.Create(ctx, order)
}

// GetOrder gets order by ID
func (s *orderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// GetOrderByCode gets order by code
func (s *orderService) GetOrderByCode(ctx context.Context, code string) (*models.Order, error) {
	return s.repo.GetByCode(ctx, code)
}

// GetOrdersByUser gets orders by user
func (s *orderService) GetOrdersByUser(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error) {
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// GetOrdersByBuyer gets orders by buyer
func (s *orderService) GetOrdersByBuyer(ctx context.Context, buyerID string, limit, offset int) ([]models.Order, int64, error) {
	return s.repo.GetByBuyerID(ctx, buyerID, limit, offset)
}

// GetOrdersByStatus gets orders by status
func (s *orderService) GetOrdersByStatus(ctx context.Context, status string, limit, offset int) ([]models.Order, int64, error) {
	return s.repo.GetByStatus(ctx, status, limit, offset)
}

// ListOrders lists all orders
func (s *orderService) ListOrders(ctx context.Context, limit, offset int) ([]models.Order, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// UpdateOrder updates order
func (s *orderService) UpdateOrder(ctx context.Context, id string, order *models.Order) error {
	return s.repo.Update(ctx, id, order)
}

// DeleteOrder deletes order
func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetOrderStats returns order statistics
func (s *orderService) GetOrderStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	// This can be expanded with actual statistics queries
	return stats, nil
}
