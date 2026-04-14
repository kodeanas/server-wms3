package services

import (
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
)

type ProductMasterService interface {
	GetByLocation(location string) ([]models.ProductMaster, error)
	GetStagingReguler() ([]dto.ProductMasterRegulerResponse, error)
	GetStagingSticker() ([]dto.ProductMasterStickerResponse, error)
	GetByBarcodeWarehouse(barcode string) (*models.ProductMaster, error)
	SetRackStaging(id string, rackStagingID string) error
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

// Get product master by barcode_warehouse
func (s *productMasterService) GetByBarcodeWarehouse(barcode string) (*models.ProductMaster, error) {
	return s.repo.FindByBarcodeWarehouse(barcode)
}

// Set rack_staging_id for product master
func (s *productMasterService) SetRackStaging(id string, rackStagingID string) error {
	return s.repo.UpdateRackStagingID(id, rackStagingID)
}
