package repositories

import (
	"context"
	"wms/models"
)

// UserRepository defines user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, int64, error)
	Update(ctx context.Context, id string, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// TaxRepository defines tax data operations
type TaxRepository interface {
	Create(ctx context.Context, tax *models.Tax) error
	GetByID(ctx context.Context, id string) (*models.Tax, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Tax, error)
	Update(ctx context.Context, id string, tax *models.Tax) error
	Delete(ctx context.Context, id string) error
}

// ProductRepository defines product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *models.ProductMaster) error
	GetByID(ctx context.Context, id string) (*models.ProductMaster, error)
	GetByBarcode(ctx context.Context, barcode string) (*models.ProductMaster, error)
	GetByCategoryID(ctx context.Context, categoryID string, limit, offset int) ([]models.ProductMaster, int64, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.ProductMaster, int64, error)
	Update(ctx context.Context, id string, product *models.ProductMaster) error
	Delete(ctx context.Context, id string) error
	GetByRackID(ctx context.Context, rackID string) ([]models.ProductMaster, error)
}

// CategoryRepository defines category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Category, int64, error)
	Update(ctx context.Context, id string, category *models.Category) error
	Delete(ctx context.Context, id string) error
}

// StickerRepository defines sticker data operations
type StickerRepository interface {
	Create(ctx context.Context, sticker *models.Sticker) error
	GetByID(ctx context.Context, id string) (*models.Sticker, error)
	GetBySlug(ctx context.Context, slug string) (*models.Sticker, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Sticker, int64, error)
	Update(ctx context.Context, id string, sticker *models.Sticker) error
	Delete(ctx context.Context, id string) error
}

// OrderRepository defines order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	GetByCode(ctx context.Context, code string) (*models.Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error)
	GetByBuyerID(ctx context.Context, buyerID string, limit, offset int) ([]models.Order, int64, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Order, int64, error)
	Update(ctx context.Context, id string, order *models.Order) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]models.Order, int64, error)
}

// CargoRepository defines cargo data operations
type CargoRepository interface {
	Create(ctx context.Context, cargo *models.Cargo) error
	GetByID(ctx context.Context, id string) (*models.Cargo, error)
	GetByCode(ctx context.Context, code string) (*models.Cargo, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Cargo, int64, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Cargo, int64, error)
	Update(ctx context.Context, id string, cargo *models.Cargo) error
	Delete(ctx context.Context, id string) error
}

// BagRepository defines bag data operations
type BagRepository interface {
	Create(ctx context.Context, bag *models.Bag) error
	GetByID(ctx context.Context, id string) (*models.Bag, error)
	GetByCode(ctx context.Context, code string) (*models.Bag, error)
	GetByCargoID(ctx context.Context, cargoID string) ([]models.Bag, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Bag, int64, error)
	Update(ctx context.Context, id string, bag *models.Bag) error
	Delete(ctx context.Context, id string) error
}

// StoreRepository defines store data operations
type StoreRepository interface {
	Create(ctx context.Context, store *models.Store) error
	GetByID(ctx context.Context, id string) (*models.Store, error)
	GetByUserID(ctx context.Context, userID string) (*models.Store, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Store, int64, error)
	Update(ctx context.Context, id string, store *models.Store) error
	Delete(ctx context.Context, id string) error
}
