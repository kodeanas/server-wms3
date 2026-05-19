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

func (s *RackDisplayService) GetAllPaginated(page, limit int, search string) ([]models.RackDisplay, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.Repo.FindAllPaginated(limit, offset, search)
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

func (s *RackDisplayService) GetSummaryAllDisplay() (*dto.RackDisplaySummaryAllResponse, error) {
	return s.Repo.GetSummaryAllDisplay()
}

// GetDetail returns rack display detail with total_item, total_price, total_price_warehouse
func (s *RackDisplayService) GetDetail(id string) (*dto.RackDisplayDetailResponse, error) {
	// Query ke product_master untuk summary
	var totalItem int
	var totalPrice float64
	var totalPriceWarehouse float64

	// Ambil data rack display
	rack, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, errors.New("Rack display not found")
	}

	db := s.Repo.DB // gunakan DB dari repository
	err = db.Table("product_masters").
		Where("rack_display_id = ? AND deleted_at IS NULL", id).
		Select("COALESCE(SUM(item),0), COALESCE(SUM(price),0), COALESCE(SUM(price_warehouse),0)").
		Row().Scan(&totalItem, &totalPrice, &totalPriceWarehouse)
	if err != nil {
		return nil, err
	}

	summary, err := s.getRackDisplaySummaryRows(id)
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
		Summary:             summary,
	}, nil
}

// getRackDisplaySummaryRows groups product_masters by category/sticker for the given rack display.
func (s *RackDisplayService) getRackDisplaySummaryRows(rackDisplayID string) ([]dto.RackDisplaySummaryItemResponse, error) {
	type row struct {
		Label          string  `gorm:"column:label"`
		Item           int64   `gorm:"column:item"`
		Price          float64 `gorm:"column:price"`
		PriceWarehouse float64 `gorm:"column:price_warehouse"`
	}

	rows := make([]row, 0)
	err := s.Repo.DB.Table("product_masters pm").
		Select(`
			CASE
				WHEN pm.category_id IS NOT NULL THEN CONCAT('category/', COALESCE(c.name, '-'))
				WHEN pm.sticker_id IS NOT NULL THEN CONCAT('sticker/', COALESCE(s.name, '-'))
				ELSE 'unknown'
			END AS label,
			COALESCE(SUM(pm.item), 0) AS item,
			COALESCE(SUM(pm.price), 0) AS price,
			COALESCE(SUM(pm.price_warehouse), 0) AS price_warehouse
		`).
		Joins("LEFT JOIN categories c ON c.id = pm.category_id::uuid").
		Joins("LEFT JOIN stickers s ON s.id = pm.sticker_id::uuid").
		Where("pm.rack_display_id = ? AND pm.deleted_at IS NULL", rackDisplayID).
		Group("label").
		Order("label ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.RackDisplaySummaryItemResponse, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.RackDisplaySummaryItemResponse{
			Label:          r.Label,
			Item:           int(r.Item),
			Price:          r.Price,
			PriceWarehouse: r.PriceWarehouse,
		})
	}
	return result, nil
}

// GetRackProductSummaryAll returns summary for all rack products across all racks.
// Contains total_item, total_price, total_price_warehouse.
func (s *RackDisplayService) GetRackProductSummaryAll() (*dto.RackProductSummaryResponse, error) {
	return s.Repo.GetRackProductSummaryAll()
}
