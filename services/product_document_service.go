package services

import (
	"wms/models"
	"wms/repositories"
)

type ProductDocumentService interface {
	ListDocuments() ([]models.ProductDocument, error)
	GetBulkDocuments() ([]models.ProductDocument, error)
	GetBulkDocumentDetail(id string) (models.ProductDocument, error)

	GetBastDocuments() ([]models.ProductDocument, error)
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
