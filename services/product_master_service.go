package services

import (
	"fmt"
	"strings"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
)

// UpdateProductMasterStagingPayload digunakan untuk update staging product master
type UpdateProductMasterStagingPayload struct {
	Price      float64 `json:"price" binding:"required"`
	Qty        int     `json:"qty" binding:"required"`
	CategoryID *string `json:"category_id"`
}

// ProductMasterService interface untuk service product master
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

// NewProductMasterService constructor
func NewProductMasterService(repo repositories.ProductMasterRepository) ProductMasterService {
	return &productMasterService{repo: repo}
}

// GetByLocation mengambil product master berdasarkan lokasi
func (s *productMasterService) GetByLocation(location string) ([]models.ProductMaster, error) {
	return s.repo.FindByLocation(location)
}

// GetStagingReguler mengambil data staging reguler
func (s *productMasterService) GetStagingReguler() ([]dto.ProductMasterRegulerResponse, error) {
	return s.repo.FindStagingReguler()
}

// GetStagingSticker mengambil data staging sticker
func (s *productMasterService) GetStagingSticker() ([]dto.ProductMasterStickerResponse, error) {
	return s.repo.FindStagingSticker()
}

// GetDetailByID mengambil detail product master berdasarkan ID
func (s *productMasterService) GetDetailByID(id string) (*dto.ProductMasterDetailResponse, error) {
	return s.repo.FindDetailByID(id)
}

// UpdateStaging melakukan update data staging product master
func (s *productMasterService) UpdateStaging(id string, input UpdateProductMasterStagingPayload) (*models.ProductMaster, error) {
	master, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Price >= 100000 {
		// Validasi category_id wajib diisi jika harga >= 100000
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
		// Jika harga < 100000, gunakan sticker
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

// stickerPriceWarehouse mengembalikan harga warehouse dari sticker
func stickerPriceWarehouse(sticker *models.Sticker) float64 {
	if sticker.FixedPrice == nil {
		return 0
	}
	return float64(*sticker.FixedPrice)
}

// GetByBarcodeWarehouse mengambil product master berdasarkan barcode warehouse
func (s *productMasterService) GetByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	return s.repo.FindByBarcodeWarehouse(barcode)
}

// SetRackStaging mengatur rack staging id pada product master
func (s *productMasterService) SetRackStaging(id string, rackStagingID string) error {
	return s.repo.UpdateRackStagingID(id, rackStagingID)
}

// ListByRackStagingID mengambil semua product master dalam satu rack staging
func (s *productMasterService) ListByRackStagingID(rackStagingID string) ([]models.ProductMaster, error) {
	return s.repo.FindAllByRackStagingID(rackStagingID)
}
