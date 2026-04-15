package services

import (
	"fmt"
	"strings"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
)

type UpdateProductMasterStagingPayload struct {
	Price      float64 `json:"price" binding:"required"`
	Qty        int     `json:"qty" binding:"required"`
	CategoryID *string `json:"category_id"`
}

type ProductMasterService interface {
	GetByLocation(location string) ([]models.ProductMaster, error)
	GetStagingReguler() ([]dto.ProductMasterRegulerResponse, error)
	GetStagingSticker() ([]dto.ProductMasterStickerResponse, error)
	GetDetailByID(id string) (*dto.ProductMasterDetailResponse, error)
	UpdateStaging(id string, input UpdateProductMasterStagingPayload) (*models.ProductMaster, error)
	GetByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	SetRackStaging(id string, rackStagingID string) error
	ListByRackStagingID(rackStagingID string) ([]models.ProductMaster, error)
}

type productMasterService struct {
	repo repositories.ProductMasterRepository
}

func NewProductMasterService(repo repositories.ProductMasterRepository) ProductMasterService {
	return &productMasterService{repo: repo}
}

func (s *productMasterService) GetByLocation(location string) ([]models.ProductMaster, error) {
	return s.repo.FindByLocation(location)
}

func (s *productMasterService) GetStagingReguler() ([]dto.ProductMasterRegulerResponse, error) {
	return s.repo.FindStagingReguler()
}

func (s *productMasterService) GetStagingSticker() ([]dto.ProductMasterStickerResponse, error) {
	return s.repo.FindStagingSticker()
}

func (s *productMasterService) GetDetailByID(id string) (*dto.ProductMasterDetailResponse, error) {
	return s.repo.FindDetailByID(id)
}

func (s *productMasterService) UpdateStaging(id string, input UpdateProductMasterStagingPayload) (*models.ProductMaster, error) {
	master, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Price >= 100000 {
		if input.CategoryID == nil || strings.TrimSpace(*input.CategoryID) == "" {
			return nil, fmt.Errorf("category_id harus diisi untuk harga >= 100000")
		}

		discount, err := s.repo.GetCategoryDiscount(*input.CategoryID)
		if err != nil {
			return nil, err
		}
		if discount == nil {
			return nil, fmt.Errorf("kategori tidak ditemukan")
		}

		master.Price = input.Price
		master.PriceWarehouse = input.Price * (1 - float64(*discount)/100)
		master.Item = input.Qty
		master.ItemWarehouse = input.Qty
		master.CategoryID = input.CategoryID
		master.StickerID = nil
		master.Location = "staging_reguler"
	} else {
		sticker, err := s.repo.FindStickerByPrice(input.Price)
		if err != nil {
			return nil, err
		}

		master.Price = input.Price
		master.PriceWarehouse = stickerPriceWarehouse(sticker)
		master.Item = input.Qty
		master.ItemWarehouse = input.Qty
		master.CategoryID = nil
		stickerID := sticker.ID.String()
		master.StickerID = &stickerID
		master.Location = "staging_sticker"
	}

	if err := s.repo.Update(master); err != nil {
		return nil, err
	}

	return master, nil
}

func stickerPriceWarehouse(sticker *models.Sticker) float64 {
	if sticker.FixedPrice == nil {
		return 0
	}
	return float64(*sticker.FixedPrice)
// Get product master by barcode_warehouse
func (s *productMasterService) GetByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	return s.repo.FindByBarcodeWarehouse(barcode)
}

// Set rack_staging_id for product master
func (s *productMasterService) SetRackStaging(id string, rackStagingID string) error {
	return s.repo.UpdateRackStagingID(id, rackStagingID)
}

// List all product master in a rack staging
func (s *productMasterService) ListByRackStagingID(rackStagingID string) ([]models.ProductMaster, error) {
	return s.repo.FindAllByRackStagingID(rackStagingID)
}
