package services

import (
	"time"
	"wms/models"
	"wms/repositories"
)

type ProductDocumentService interface {
	ListDocuments() ([]models.ProductDocument, error)
	GetBulkDocuments() ([]models.ProductDocument, error)
	GetBulkDocumentDetail(id string) (models.ProductDocument, error)

	GetBastDocuments() ([]models.ProductDocument, error)
	FinishDocument(id string) error
}

type productDocumentService struct {
	repo repositories.ProductDocumentRepository
}

func NewProductDocumentService(repo repositories.ProductDocumentRepository) ProductDocumentService {
	return &productDocumentService{repo: repo}
}

func (s *productDocumentService) ListDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindAll()
}

func (s *productDocumentService) FinishDocument(id string) error {
	now := time.Now()
	err := s.repo.UpdateDateStopByID(id, &now)
	if err != nil {
		return err
	}
	return s.repo.UpdateStatusByID(id, "done")
}

// Implementasi filter bulk
func (s *productDocumentService) GetBulkDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindByType("bulk")
}

func (s *productDocumentService) GetBulkDocumentDetail(id string) (models.ProductDocument, error) {
	return s.repo.FindBulkDetailByID(id)
}

// Implementasi filter bast
func (s *productDocumentService) GetBastDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindByType("bast")
}
