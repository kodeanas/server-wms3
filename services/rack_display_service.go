package services

import (
	"errors"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
)

type RackDisplayService struct {
	Repo *repositories.RackDisplayRepository
}

func NewRackDisplayService(repo *repositories.RackDisplayRepository) *RackDisplayService {
	return &RackDisplayService{Repo: repo}
}

func (s *RackDisplayService) Create(rack *models.RackDisplay) error {
	return s.Repo.Create(rack)
}

func (s *RackDisplayService) GetAll() ([]models.RackDisplay, error) {
	return s.Repo.FindAll()
}

func (s *RackDisplayService) GetByID(id string) (*models.RackDisplay, error) {
	return s.Repo.FindByID(id)
}

func (s *RackDisplayService) Update(rack *models.RackDisplay) error {
	return s.Repo.Update(rack)
}

func (s *RackDisplayService) Delete(id string) error {
	return s.Repo.SoftDelete(id)
}

// GetDetail returns rack display detail with total_item, total_price, total_price_warehouse
func (s *RackDisplayService) GetDetail(id string) (*dto.RackDisplayDetailResponse, error) {
	// Ambil data rack display
	rack, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, errors.New("Rack display not found")
	}

	// Query ke product_master untuk summary
	var totalItem int
	var totalPrice float64
	var totalPriceWarehouse float64

	db := s.Repo.DB // gunakan DB dari repository
	err = db.Table("product_masters").
		Where("rack_display_id = ? AND deleted_at IS NULL", id).
		Select("COALESCE(SUM(item),0), COALESCE(SUM(price),0), COALESCE(SUM(price_warehouse),0)").
		Row().Scan(&totalItem, &totalPrice, &totalPriceWarehouse)
	if err != nil {
		return nil, err
	}

	return &dto.RackDisplayDetailResponse{
		ID:                  rack.ID.String(),
		Code:                rack.Code,
		Name:                rack.Name,
		CreatedAt:           rack.CreatedAt,
		TotalItem:           totalItem,
		TotalPrice:          totalPrice,
		TotalPriceWarehouse: totalPriceWarehouse,
	}, nil
}
