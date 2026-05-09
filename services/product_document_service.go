package services

import (
	"encoding/json"
	"math"
	"time"
	"wms/models"
	"wms/repositories"
)

// ProductDocumentService interface untuk service dokumen produk
type ProductDocumentService interface {
	ListDocuments() ([]models.ProductDocument, error)
	ListDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error)
	GetBulkDocuments() ([]models.ProductDocument, error)
	GetBulkDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error)
	GetBulkDocumentDetail(id string) (models.ProductDocument, error)
	GetBastDocuments() ([]models.ProductDocument, error)
	GetBastDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error)
	GetBastRelationsDetail(id string) (map[string]interface{}, error)
	GetBastOverview(id string) (map[string]interface{}, error)
	GetBastPendingsByType(id string) (map[string]interface{}, error)
	FinishDocument(id string) error
}

type productDocumentService struct {
	repo repositories.ProductDocumentRepository
}

// NewProductDocumentService constructor
func NewProductDocumentService(repo repositories.ProductDocumentRepository) ProductDocumentService {
	return &productDocumentService{repo: repo}
}

// ListDocuments mengambil semua dokumen produk
func (s *productDocumentService) ListDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindAll()
}

func (s *productDocumentService) ListDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindAllPaginated(limit, offset, search)
}

// FinishDocument menandai dokumen selesai
func (s *productDocumentService) FinishDocument(id string) error {
	now := time.Now()
	if err := s.repo.UpdateDateStopByID(id, &now); err != nil {
		return err
	}
	return s.repo.UpdateStatusByID(id, "done")
}

// GetBulkDocuments mengambil dokumen dengan tipe bulk
func (s *productDocumentService) GetBulkDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindByType("bulk")
}

func (s *productDocumentService) GetBulkDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindByTypePaginated("bulk", limit, offset, search)
}

// GetBulkDocumentDetail mengambil detail dokumen bulk berdasarkan ID
func (s *productDocumentService) GetBulkDocumentDetail(id string) (models.ProductDocument, error) {
	return s.repo.FindBulkDetailByID(id)
}

// GetBastDocuments mengambil dokumen dengan tipe bast
func (s *productDocumentService) GetBastDocuments() ([]models.ProductDocument, error) {
	return s.repo.FindByType("bast")
}

func (s *productDocumentService) GetBastDocumentsPaginated(page, limit int, search string) ([]models.ProductDocument, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.FindByTypePaginated("bast", limit, offset, search)
}

// GetBastRelationsDetail mengambil detail relasi bast (pending & scanned)
func (s *productDocumentService) GetBastRelationsDetail(id string) (map[string]interface{}, error) {
	if _, err := s.repo.FindBastByID(id); err != nil {
		return nil, err
	}

	productPending, err := s.repo.FindBastProductPendingByDiscrepancy(id)
	if err != nil {
		return nil, err
	}

	productScanned, err := s.repo.FindBastProductPendingByNonDiscrepancy(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"product_pending": productPending,
		"product_scanned": productScanned,
	}, nil
}

// GetBastOverview mengambil overview dokumen bast
func (s *productDocumentService) GetBastOverview(id string) (map[string]interface{}, error) {
	doc, err := s.repo.FindBastByID(id)
	if err != nil {
		return nil, err
	}

	totalItemScanned, totalPriceScanned, err := s.repo.FindBastScannedSummary(id)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(docBytes, &data); err != nil {
		return nil, err
	}

	data["total_item_scanned"] = totalItemScanned
	data["total_price_scanned"] = totalPriceScanned

	return data, nil
}

// GetBastPendingsByType mengambil summary pending bast berdasarkan status
func (s *productDocumentService) GetBastPendingsByType(id string) (map[string]interface{}, error) {
	doc, err := s.repo.FindBastByID(id)
	if err != nil {
		return nil, err
	}

	statuses := []string{"good", "damaged", "abnormal", "non", "discrepancy"}
	summaryByStatus, err := s.repo.FindBastPendingSummaryByStatuses(id, statuses)
	if err != nil {
		return nil, err
	}

	totalScannedItem, totalScannedPrice, err := s.repo.FindBastScannedSummary(id)
	if err != nil {
		return nil, err
	}

	totalItemInFile := doc.FileItem
	totalItemInFileFloat := float64(totalItemInFile)
	totalPriceInFileFloat := float64(doc.FilePrice)

	round2 := func(v float64) float64 {
		return math.Round(v*100) / 100
	}

	buildSummary := func(status string) map[string]float64 {
		statusData := summaryByStatus[status]
		totalItem := statusData["total_item"]
		totalPrice := statusData["total_price"]

		percentageItem := 0.0
		percentagePrice := 0.0
		if totalItemInFileFloat > 0 {
			percentageItem = round2((totalItem / totalItemInFileFloat) * 100)
		}
		if totalPriceInFileFloat > 0 {
			percentagePrice = round2((totalPrice / totalPriceInFileFloat) * 100)
		}

		return map[string]float64{
			"total_item":       totalItem,
			"total_price":      totalPrice,
			"percentage_item":  percentageItem,
			"percentage_price": percentagePrice,
		}
	}

	totalScannedSummary := map[string]float64{
		"total_item":       float64(totalScannedItem),
		"total_price":      totalScannedPrice,
		"percentage_item":  0,
		"percentage_price": 0,
	}
	if totalItemInFileFloat > 0 {
		totalScannedSummary["percentage_item"] = round2((float64(totalScannedItem) / totalItemInFileFloat) * 100)
	}
	if totalPriceInFileFloat > 0 {
		totalScannedSummary["percentage_price"] = round2((totalScannedPrice / totalPriceInFileFloat) * 100)
	}

	totalItemInFileSummary := map[string]float64{
		"total_item":       totalItemInFileFloat,
		"total_price":      totalPriceInFileFloat,
		"percentage_item":  100,
		"percentage_price": 100,
	}

	return map[string]interface{}{
		"good":               buildSummary("good"),
		"damaged":            buildSummary("damaged"),
		"abnormal":           buildSummary("abnormal"),
		"non":                buildSummary("non"),
		"total_discrepacy":   buildSummary("discrepancy"),
		"total_scanned":      totalScannedSummary,
		"total_item_in_file": totalItemInFileSummary,
	}, nil
}
