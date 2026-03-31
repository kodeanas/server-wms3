package repositories

import (
	"context"
	"wms/config"
	"wms/models"
)

type productRepository struct{}

// NewProductRepository creates a new product repository
func NewProductRepository() ProductRepository {
	return &productRepository{}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *models.ProductMaster) error {
	return config.DB.WithContext(ctx).Create(product).Error
}

// GetByID gets product by ID
func (r *productRepository) GetByID(ctx context.Context, id string) (*models.ProductMaster, error) {
	var product models.ProductMaster
	err := config.DB.WithContext(ctx).
		Preload("Category").
		Preload("Sticker").
		Preload("User").
		First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetByBarcode gets product by barcode
func (r *productRepository) GetByBarcode(ctx context.Context, barcode string) (*models.ProductMaster, error) {
	var product models.ProductMaster
	err := config.DB.WithContext(ctx).
		Preload("Category").
		Preload("Sticker").
		Preload("User").
		First(&product, "barcode = ?", barcode).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetByCategoryID gets products by category
func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID string, limit, offset int) ([]models.ProductMaster, int64, error) {
	var products []models.ProductMaster
	var total int64

	err := config.DB.WithContext(ctx).
		Model(&models.ProductMaster{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("Category").
		Preload("Sticker").
		Preload("User").
		Limit(limit).Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetAll gets all products
func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]models.ProductMaster, int64, error) {
	var products []models.ProductMaster
	var total int64

	err := config.DB.WithContext(ctx).Model(&models.ProductMaster{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = config.DB.WithContext(ctx).
		Preload("Category").
		Preload("Sticker").
		Preload("User").
		Limit(limit).Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update updates product
func (r *productRepository) Update(ctx context.Context, id string, product *models.ProductMaster) error {
	return config.DB.WithContext(ctx).Model(&models.ProductMaster{}).Where("id = ?", id).Updates(product).Error
}

// Delete deletes product
func (r *productRepository) Delete(ctx context.Context, id string) error {
	return config.DB.WithContext(ctx).Delete(&models.ProductMaster{}, "id = ?", id).Error
}

// GetByRackID gets products by rack
func (r *productRepository) GetByRackID(ctx context.Context, rackID string) ([]models.ProductMaster, error) {
	var products []models.ProductMaster
	err := config.DB.WithContext(ctx).
		Where("rack_id = ?", rackID).
		Preload("Category").
		Preload("Sticker").
		Preload("User").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
