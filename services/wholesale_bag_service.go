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
	ListWholesaleBags() ([]models.Bag, error)
	GetWholesaleBagByID(id string) (*models.Bag, error)
	ListProductsByWholesaleBagID(bagID string) ([]models.ProductMaster, error)
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
	var uidPtr *uuid.UUID
	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err == nil {
			uidPtr = &uid
		}
	}
	bag := &models.Bag{
		Code:      fmt.Sprintf("WHOLESALE-%d", time.Now().UnixNano()),
		Type:      "reguler",
		IsMoved:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if uidPtr != nil {
		bag.UserID = *uidPtr
	}
	if err := s.repo.Create(bag); err != nil {
		return nil, err
	}
	return bag, nil
}

func (s *wholesaleBagService) ListWholesaleBags() ([]models.Bag, error) {
       return s.repo.FindByType("reguler")
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
