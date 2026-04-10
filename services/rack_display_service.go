package services

import (
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
