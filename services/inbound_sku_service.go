package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"wms/models"
	"wms/repositories"
	"wms/utils"

	"github.com/google/uuid"
)

type InboundSKUService interface {
	UploadExcelAndCreatePendings(file interface{}, fileType, supplier, fileName string) (int, int, []string, error)
	CrosscheckPending(pendingID string, itemGood, itemDamaged int) error
	FinishInboundSKU(documentID string) error
	ListSKUProductDocuments() ([]models.ProductDocument, error)
	// List all SKU product documents
	GetPendingByID(pendingID string) (*models.ProductPending, error)
	GetDocumentByID(documentID string) (*models.ProductDocument, error)
}

type inboundSKUService struct {
	productDocumentRepo repositories.ProductDocumentRepository
	productPendingRepo  repositories.ProductPendingRepository
	productRepairRepo   repositories.ProductRepairRepository
	productMasterRepo   repositories.ProductMasterRepository
}

func NewInboundSKUService(
	productDocumentRepo repositories.ProductDocumentRepository,
	productPendingRepo repositories.ProductPendingRepository,
	productRepairRepo repositories.ProductRepairRepository,
	productMasterRepo repositories.ProductMasterRepository,
) InboundSKUService {
	return &inboundSKUService{
		productDocumentRepo: productDocumentRepo,
		productPendingRepo:  productPendingRepo,
		productRepairRepo:   productRepairRepo,
		productMasterRepo:   productMasterRepo,
	}
}

func (s *inboundSKUService) UploadExcelAndCreatePendings(file interface{}, fileType, supplier, fileName string) (int, int, []string, error) {
	mf, ok := file.(multipart.File)
	if !ok {
		return 0, 0, nil, fmt.Errorf("invalid file type, must be multipart.File")
	}
	defer mf.Close()
	_, rows, err := utils.ParseBulkFile(mf, fileType)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("parse error: %w", err)
	}
	if len(rows) == 0 {
		return 0, 0, nil, fmt.Errorf("no data found in file")
	}

	idxItem := 4  // index kolom item
	idxPrice := 3 // index kolom harga

	fileItem := 0
	filePrice := 0
	for _, row := range rows {
		if len(row) <= idxItem || len(row) <= idxPrice {
			continue
		}
		item := parseInt(row[idxItem])
		price := parseFloat(row[idxPrice])
		fileItem += item
		filePrice += int(price * float64(item))
	}

	doc := &models.ProductDocument{
		ID:        uuid.New(),
		Code:      uuid.New().String(),
		FileName:  fileName,
		FileItem:  fileItem,
		FilePrice: filePrice,
		Status:    "pending",
		Type:      "sku",
		Supplier:  supplier,
	}
	err = s.productDocumentRepo.Create(doc)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to create document: %w", err)
	}

	inserted := 0
	skipped := 0
	skipDetails := []string{}
	for idx, row := range rows {
		if len(row) < 4 {
			skipped++
			skipDetails = append(skipDetails, fmt.Sprintf("Row %d skipped: kolom kurang lengkap: %v", idx+1, row))
			continue
		}
		pending := &models.ProductPending{
			ID:         uuid.New(),
			DocumentID: doc.ID.String(),
			Barcode:    row[0],
			Name:       row[1],
			Price:      parseFloat(row[3]),
			Item:       parseInt(row[4]),
			IsSKU:      true,
			Status:     "good",
		}
		err = s.productPendingRepo.Create(pending)
		if err != nil {
			skipped++
			skipDetails = append(skipDetails, fmt.Sprintf("Row %d skipped: gagal insert product_pending: %v | DB error: %v", idx+1, row, err))
			continue
		}
		inserted++
	}
	return inserted, skipped, skipDetails, nil
}

func (s *inboundSKUService) ListSKUProductDocuments() ([]models.ProductDocument, error) {
	return s.productDocumentRepo.FindByType("sku")
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func (s *inboundSKUService) CrosscheckPending(pendingID string, itemGood, itemDamaged int) error {
	pending, err := s.productPendingRepo.FindByID(pendingID)
	if err != nil {
		return fmt.Errorf("pending not found: %w", err)
	}
	if itemGood+itemDamaged > pending.Item {
		return fmt.Errorf("item_good + item_damaged tidak boleh lebih dari total item")
	}
	pending.ItemGood = itemGood
	pending.ItemDamaged = itemDamaged
	return s.productPendingRepo.Update(pending)
}

// Helper untuk konversi string ke *string
func toStrPtr(s string) *string {
	return &s
}

func (s *inboundSKUService) FinishInboundSKU(documentID string) error {
	// Cek dokumen sudah selesai
	doc, err := s.productDocumentRepo.FindSkuDetailByID(documentID)
	if err != nil {
		return fmt.Errorf("document not found or not SKU type: %w", err)
	}
	if doc.Status == "done" {
		return fmt.Errorf("document already finished (done)")
	}
	// Cek dokumen ada dan type sku
	_, err = s.productDocumentRepo.FindSkuDetailByID(documentID)
	if err != nil {
		return fmt.Errorf("document not found or not SKU type: %w", err)
	}

	pendings, err := s.productPendingRepo.FindByDocumentID(documentID)
	if err != nil {
		return err
	}
	for _, pending := range pendings {
		itemGood := pending.ItemGood
		itemDamaged := pending.ItemDamaged
		fmt.Printf("DEBUG: pending_id=%s, barcode=%s, item_good=%d, item_damaged=%d\n", pending.ID.String(), pending.Barcode, itemGood, itemDamaged)

		// Proses ke product_master
		if itemGood > 0 {
			master := &models.ProductMaster{
				ID:               uuid.New(),
				DocumentID:       pending.DocumentID,
				Barcode:          pending.Barcode,
				Name:             pending.Name,
				ItemWarehouse:    itemGood,
				PriceWarehouse:   pending.Price,
				IsSKU:            true,
				Location:         "staging_sku",
				ProductPendingID: toStrPtr(pending.ID.String()),
			}
			err := s.productMasterRepo.Create(master)
			if err != nil {
				fmt.Printf("DEBUG: Gagal insert ke product_master: %v\n", err)
			} else {
				fmt.Printf("DEBUG: Berhasil insert ke product_master: barcode=%s, item_warehouse=%d\n", master.Barcode, master.ItemWarehouse)
			}
		}
		// Proses ke product_repair
		if itemDamaged > 0 {
			repair := &models.ProductRepair{
				ID:          uuid.New(),
				ProductID:   &pending.ID,
				Status:      "progress",
				ItemBefore:  pending.Item,
				ItemUpdate:  itemDamaged,
				PriceBefore: pending.Price,
				PriceUpdate: pending.Price,
			}
			err := s.productRepairRepo.Create(repair)
			if err != nil {
				fmt.Printf("DEBUG: Gagal insert ke product_repair: %v\n", err)
			} else {
				fmt.Printf("DEBUG: Berhasil insert ke product_repair: item_update=%d\n", itemDamaged)
			}
		}
	}
	_ = s.productDocumentRepo.UpdateStatusByID(documentID, "done")
	return nil
}

// Get pending by ID
func (s *inboundSKUService) GetPendingByID(pendingID string) (*models.ProductPending, error) {
	return s.productPendingRepo.FindByID(pendingID)
}

// Get document by ID
func (s *inboundSKUService) GetDocumentByID(documentID string) (*models.ProductDocument, error) {
	doc, err := s.productDocumentRepo.FindSkuDetailByID(documentID)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
