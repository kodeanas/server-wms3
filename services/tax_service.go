package services

import (
	"wms/models"
	"wms/repositories"

	"github.com/google/uuid"
)

type TaxService struct {
	repo *repositories.TaxRepository
}

func NewTaxService(repo *repositories.TaxRepository) *TaxService {
	return &TaxService{repo: repo}
}

func (s *TaxService) Create(tax *models.Tax) error {
	if tax.ID == (models.Tax{}).ID {
		// Generate UUID jika belum ada
		tax.ID = uuid.New()
	}
	err := s.repo.Create(tax)
	if err != nil {
		return err
	}
	if tax.IsActive {
		// Set all other taxes inactive setelah create
		s.repo.SetAllInactiveExcept(tax.ID.String())
	}
	return nil
}

func (s *TaxService) Update(tax *models.Tax) error {
	if tax.IsActive {
		s.repo.SetAllInactiveExcept(tax.ID.String())
	}
	return s.repo.Update(tax)
}

func (s *TaxService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *TaxService) FindByID(id string) (*models.Tax, error) {
	return s.repo.FindByID(id)
}

func (s *TaxService) FindAllPaginated(page, limit int, search string) ([]models.Tax, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAllPaginated(limit, offset, search)
}

func (s *TaxService) FindActive() (*models.Tax, error) {
	return s.repo.FindActive()
}
