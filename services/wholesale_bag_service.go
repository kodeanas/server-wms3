package services

import (
	"fmt"
	"time"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

type WholesaleBagService interface {
	CreateWholesaleBag(userID string) (*models.Bag, error)
	ListWholesaleBags() ([]dto.WholesaleBagListResponse, error)
	ListWholesaleBagsPaginated(page, limit int, search string) ([]dto.WholesaleBagListResponse, int64, error)
	GetWholesaleBagByID(id string) (*models.Bag, error)
	ListProductsByWholesaleBagID(bagID string) ([]models.ProductMaster, error)
	ListProductsByWholesaleBagIDPaginated(bagID string, page, limit int, search string) ([]models.ProductMaster, int64, error)
	GetWholesaleBagDetail(bagID string) (*dto.RackStagingDetailResponse, error)
	GetProductByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	SetBag(productID string, bagID string) error
}

type wholesaleBagService struct {
	repo              repositories.BagRepository
	productMasterRepo repositories.ProductMasterRepository
}

func NewWholesaleBagService(repo repositories.BagRepository, productMasterRepo repositories.ProductMasterRepository) WholesaleBagService {
	return &wholesaleBagService{repo: repo, productMasterRepo: productMasterRepo}
}

func (s *wholesaleBagService) CreateWholesaleBag(userID string) (*models.Bag, error) {
	bag := &models.Bag{
		Code:      fmt.Sprintf("WHOLESALE-%d", time.Now().UnixNano()),
		Type:      "reguler",
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

func (s *wholesaleBagService) ListWholesaleBags() ([]dto.WholesaleBagListResponse, error) {
	bags, err := s.repo.FindByType("reguler")
	if err != nil {
		return nil, err
	}
	var resp []dto.WholesaleBagListResponse
	for _, b := range bags {
		resp = append(resp, dto.WholesaleBagListResponse{
			ID:        b.ID.String(),
			Code:      b.Code,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
			DeletedAt: b.DeletedAt,
		})
	}
	return resp, nil
}

func (s *wholesaleBagService) ListWholesaleBagsPaginated(page, limit int, search string) ([]dto.WholesaleBagListResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	bags, total, err := s.repo.FindByTypePaginated("reguler", limit, offset, search)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]dto.WholesaleBagListResponse, 0, len(bags))
	for _, b := range bags {
		resp = append(resp, dto.WholesaleBagListResponse{
			ID:        b.ID.String(),
			Code:      b.Code,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
			DeletedAt: b.DeletedAt,
		})
	}
	return resp, total, nil
}

func (s *wholesaleBagService) GetWholesaleBagByID(id string) (*models.Bag, error) {
	bag, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if bag.Type != "reguler" {
		return nil, fmt.Errorf("Bag bukan tipe reguler")
	}
	return bag, nil
}

func (s *wholesaleBagService) ListProductsByWholesaleBagID(bagID string) ([]models.ProductMaster, error) {
	return s.productMasterRepo.FindByBagID(bagID)
}

func (s *wholesaleBagService) ListProductsByWholesaleBagIDPaginated(bagID string, page, limit int, search string) ([]models.ProductMaster, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.productMasterRepo.FindByBagIDPaginated(bagID, limit, offset, search)
}

func (s *wholesaleBagService) GetWholesaleBagDetail(bagID string) (*dto.RackStagingDetailResponse, error) {
	bag, err := s.repo.FindByID(bagID)
	if err != nil {
		return nil, err
	}
	if bag.Type != "reguler" {
		return nil, fmt.Errorf("Bag bukan tipe reguler")
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
		RackDisplayName:     "-",
		CreatedAt:           bag.CreatedAt.Format(time.RFC3339),
		IsMoved:             bag.IsMoved,
		TotalItem:           totalItem,
		TotalPriceWarehouse: totalPrice,
	}
	return resp, nil
}

func (s *wholesaleBagService) GetProductByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	return s.productMasterRepo.FindByBarcodeWarehouse(barcode)
}

func (s *wholesaleBagService) SetBag(productID string, bagID string) error {
	return s.productMasterRepo.UpdateBagID(productID, bagID)
}
