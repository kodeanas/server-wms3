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
