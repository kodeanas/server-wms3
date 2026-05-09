package services

import (
	"fmt"
	"time"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

type RackStagingStickerService interface {
	CreateStickerBag(userID string) (*models.Bag, error)
	GetBagByID(id string) (*models.Bag, error)
	ListBags() ([]models.Bag, error)
	ListBagsPaginated(page, limit int, search string) ([]models.Bag, int64, error)
	GetProductByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	SetBag(productID string, bagID string) error
	ListProductsByBagID(bagID string) ([]models.ProductMaster, error)
	ListProductsByBagIDPaginated(bagID string, page, limit int, search string) ([]models.ProductMaster, int64, error)
	GetBagDetail(bagID string) (*dto.RackStagingDetailResponse, error)
}

func (s *rackStagingStickerService) ListProductsByBagID(bagID string) ([]models.ProductMaster, error) {
	return s.productMasterRepo.FindByBagID(bagID)
}

func (s *rackStagingStickerService) ListProductsByBagIDPaginated(bagID string, page, limit int, search string) ([]models.ProductMaster, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.productMasterRepo.FindByBagIDPaginated(bagID, limit, offset, search)
}

type rackStagingStickerService struct {
	repo              repositories.BagRepository
	productMasterRepo repositories.ProductMasterRepository
}

func NewRackStagingStickerService(repo repositories.BagRepository, productMasterRepo repositories.ProductMasterRepository) RackStagingStickerService {
	return &rackStagingStickerService{repo: repo, productMasterRepo: productMasterRepo}
}

func generateBagCode() string {
	return fmt.Sprintf("BAG-%d", time.Now().UnixNano())
}

func (s *rackStagingStickerService) CreateStickerBag(userID string) (*models.Bag, error) {
	bag := &models.Bag{
		Code:      generateBagCode(),
		Type:      "sticker",
		IsMoved:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil && uid != uuid.Nil {
			bag.UserID = &uid
		}
	}
	if err := s.repo.Create(bag); err != nil {
		return nil, err
	}
	return bag, nil
}

func (s *rackStagingStickerService) GetBagByID(id string) (*models.Bag, error) {
	return s.repo.FindByID(id)
}

func (s *rackStagingStickerService) ListBags() ([]models.Bag, error) {
	return s.repo.FindByType("sticker")
}

func (s *rackStagingStickerService) ListBagsPaginated(page, limit int, search string) ([]models.Bag, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindByTypePaginated("sticker", limit, offset, search)
}

func (s *rackStagingStickerService) GetProductByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	return s.productMasterRepo.FindByBarcodeWarehouse(barcode)
}

func (s *rackStagingStickerService) SetBag(productID string, bagID string) error {
	return s.productMasterRepo.UpdateBagID(productID, bagID)
}

func (s *rackStagingStickerService) GetBagDetail(bagID string) (*dto.RackStagingDetailResponse, error) {
	bag, err := s.repo.FindByID(bagID)
	if err != nil {
		return nil, err
	}
	products, err := s.productMasterRepo.FindByBagID(bagID)
	if err != nil {
		return nil, err
	}
	totalItem := 0
	totalPrice := 0.0
	for _, pm := range products {
		totalItem += pm.ItemWarehouse
		totalPrice += pm.PriceWarehouse
	}
	resp := &dto.RackStagingDetailResponse{
		Code:                bag.Code,
		CreatedAt:           bag.CreatedAt.Format(time.RFC3339),
		IsMoved:             bag.IsMoved,
		TotalItem:           totalItem,
		TotalPriceWarehouse: totalPrice,
	}
	return resp, nil
}
