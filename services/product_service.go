package services

import (
	"context"
	"wms/models"
	"wms/repositories"
)

// ProductService defines product business logic
type ProductService interface {
	CreateProduct(ctx context.Context, product *models.ProductMaster) error
	GetProduct(ctx context.Context, id string) (*models.ProductMaster, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*models.ProductMaster, error)
	GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) ([]models.ProductMaster, int64, error)
	ListProducts(ctx context.Context, limit, offset int) ([]models.ProductMaster, int64, error)
	GetProductsByRack(ctx context.Context, rackID string) ([]models.ProductMaster, error)
	UpdateProduct(ctx context.Context, id string, product *models.ProductMaster) error
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo repositories.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, product *models.ProductMaster) error {
	return s.repo.Create(ctx, product)
}

// GetProduct gets product by ID
func (s *productService) GetProduct(ctx context.Context, id string) (*models.ProductMaster, error) {
	return s.repo.GetByID(ctx, id)
}

// GetProductByBarcode gets product by barcode
func (s *productService) GetProductByBarcode(ctx context.Context, barcode string) (*models.ProductMaster, error) {
	return s.repo.GetByBarcode(ctx, barcode)
}

// GetProductsByCategory gets products by category
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) ([]models.ProductMaster, int64, error) {
	return s.repo.GetByCategoryID(ctx, categoryID, limit, offset)
}

// ListProducts lists all products
func (s *productService) ListProducts(ctx context.Context, limit, offset int) ([]models.ProductMaster, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// GetProductsByRack gets products by rack
func (s *productService) GetProductsByRack(ctx context.Context, rackID string) ([]models.ProductMaster, error) {
	return s.repo.GetByRackID(ctx, rackID)
}

// UpdateProduct updates product
func (s *productService) UpdateProduct(ctx context.Context, id string, product *models.ProductMaster) error {
	return s.repo.Update(ctx, id, product)
}

// DeleteProduct deletes product
func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
