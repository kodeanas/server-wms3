package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"wms/models"
	"wms/repositories"
	"wms/utils"

	"github.com/google/uuid"
)

type InboundSKUService interface {
	UploadExcelAndCreatePendings(file interface{}, fileType, supplier, fileName string) error
	CrosscheckPending(pendingID string, itemGood, itemDamaged int) error
	FinishInboundSKU(documentID string) error
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

func (s *inboundSKUService) UploadExcelAndCreatePendings(file interface{}, fileType, supplier, fileName string) error {
	mf, ok := file.(multipart.File)
	if !ok {
		return fmt.Errorf("invalid file type, must be multipart.File")
	}
	defer mf.Close()
	_, rows, err := utils.ParseBulkFile(mf, fileType)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("no data found in file")
	}

	doc := &models.ProductDocument{
		ID:       uuid.New(),
		Code:     uuid.New().String(),
		FileName: fileName,
		Status:   "pending",
		Type:     "sku",
		Supplier: supplier,
	}
	err = s.productDocumentRepo.Create(doc)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	for _, row := range rows {
		if len(row) < 4 {
			continue // skip baris tidak valid
		}
		pending := &models.ProductPending{
			ID:         uuid.New(),
			DocumentID: doc.ID.String(),
			Barcode:    row[0],
			Name:       row[1],
			Price:      parseFloat(row[2]),
			Item:       parseInt(row[3]),
			IsSKU:      true,
			Status:     "pending",
		}
		err = s.productPendingRepo.Create(pending)
		if err != nil {
			return fmt.Errorf("failed to create pending: %w", err)
		}
	}
	return nil
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
	pending.Note = fmt.Sprintf("good: %d, damaged: %d", itemGood, itemDamaged)
	pending.Status = "crosschecked"
	// Simpan nilai item_good dan item_damaged di Note, atau tambahkan field baru jika perlu
	return s.productPendingRepo.Update(pending)
}

func (s *inboundSKUService) FinishInboundSKU(documentID string) error {
	pendings, err := s.productPendingRepo.FindByDocumentID(documentID)
	if err != nil {
		return err
	}
	for _, pending := range pendings {
		// Ambil item_good dan item_damaged dari Note
		itemGood, itemDamaged := 0, 0
		if strings.Contains(pending.Note, "good:") && strings.Contains(pending.Note, "damaged:") {
			_, err := fmt.Sscanf(pending.Note, "good: %d, damaged: %d", &itemGood, &itemDamaged)
			if err != nil {
				itemGood, itemDamaged = 0, 0
			}
		}

		// Proses ke product_master
		if itemGood > 0 {
			// TODO: Cek apakah sudah ada di product_master, jika ada update, jika tidak insert
			// Sementara: insert baru
			master := &models.ProductMaster{
				ID:             uuid.New(),
				DocumentID:     pending.DocumentID,
				Barcode:        pending.Barcode,
				Name:           pending.Name,
				ItemWarehouse:  itemGood,
				PriceWarehouse: pending.Price,
				IsSKU:          true,
				Location:       "warehouse",
			}
			_ = s.productMasterRepo.Create(master) // TODO: handle error
		}

		// Proses ke product_repair
		if itemDamaged > 0 {
			repair := &models.ProductRepair{
				ID:          uuid.New(),
				DocumentID:  pending.DocumentID,
				Barcode:     pending.Barcode,
				Name:        pending.Name,
				ItemBefore:  pending.Item,
				ItemUpdate:  itemDamaged,
				PriceBefore: pending.Price,
				PriceUpdate: pending.Price,
			}
			_ = s.productRepairRepo.Create(repair) // TODO: handle error
		}

		// Item hilang tidak diproses
	}
	_ = s.productDocumentRepo.UpdateStatusByID(documentID, "finished")
	return nil
}
