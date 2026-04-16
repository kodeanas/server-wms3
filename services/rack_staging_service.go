package services

import (
	"fmt"
	"time"

	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
)

type RackStagingService struct {
	RackStagingRepo *repositories.RackStagingRepository
	RackDisplayRepo *repositories.RackDisplayRepository
}

func NewRackStagingService(rackStagingRepo *repositories.RackStagingRepository, rackDisplayRepo *repositories.RackDisplayRepository) *RackStagingService {
	return &RackStagingService{
		RackStagingRepo: rackStagingRepo,
		RackDisplayRepo: rackDisplayRepo,
	}
}

func (s *RackStagingService) CreateRackStaging(rackDisplayID string) (*dto.RackStagingResponse, error) {
	display, err := s.RackDisplayRepo.FindByID(rackDisplayID)
	if err != nil {
		return nil, fmt.Errorf("rack display not found: %w", err)
	}

	count, err := s.RackStagingRepo.CountByRackDisplayID(rackDisplayID)
	if err != nil {
		return nil, err
	}
	order := count + 1

	name := fmt.Sprintf("%s - %d", display.Name, order)
	code := fmt.Sprintf("RSTG-%s-%d", display.Code, order)

	rack := &models.RackStaging{
		RackDisplayID: display.ID,
		Code:          code,
		Name:          name,
		IsMoved:       false,
	}

	err = s.RackStagingRepo.Create(rack)
	if err != nil {
		return nil, err
	}

	resp := dto.RackStagingResponse{
		ID:            rack.ID.String(),
		RackDisplayID: rack.RackDisplayID.String(),
		Code:          rack.Code,
		Name:          rack.Name,
		IsMoved:       rack.IsMoved,
		CreatedAt:     rack.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     rack.UpdatedAt.Format(time.RFC3339),
	}
	return &resp, nil
}

// Get detail of a rack staging by ID
func (s *RackStagingService) GetRackStagingDetail(id string) (*dto.RackStagingDetailResponse, error) {
	rack, err := s.RackStagingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	display, err := s.RackDisplayRepo.FindByID(rack.RackDisplayID.String())
	if err != nil {
		return nil, err
	}
	// Ambil product master yang ada di rack staging ini
	productMasters, err := repositories.NewProductMasterRepository(s.RackStagingRepo.DB).FindAllByRackStagingID(id)
	totalItem := 0
	totalPrice := 0.0
	if err == nil {
		totalItem = len(productMasters)
		for _, pm := range productMasters {
			totalPrice += pm.PriceWarehouse
		}
	}
	resp := &dto.RackStagingDetailResponse{
		Code:                rack.Code,
		RackDisplayName:     display.Name,
		CreatedAt:           rack.CreatedAt.Format(time.RFC3339),
		IsMoved:             rack.IsMoved,
		TotalItem:           totalItem,
		TotalPriceWarehouse: totalPrice,
	}
	return resp, nil
}

// List all rack stagings
func (s *RackStagingService) ListAllRackStaging() ([]models.RackStaging, error) {
	return s.RackStagingRepo.FindAllRackStaging()
}

// Finish rack staging: set is_moved, update semua product master ke display
func (s *RackStagingService) FinishRackStaging(rackStagingID string, productMasterRepo repositories.ProductMasterRepository) error {
	rack, err := s.RackStagingRepo.FindByID(rackStagingID)
	if err != nil {
		return err
	}
	// Set is_moved = true
	if err := s.RackStagingRepo.SetIsMoved(rackStagingID); err != nil {
		return err
	}
	// Update semua product master di rack staging ini
	if err := productMasterRepo.MoveAllToDisplay(rackStagingID, rack.RackDisplayID.String()); err != nil {
		return err
	}
	return nil
}
