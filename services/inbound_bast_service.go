package services

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	response "wms/dto/response"
	"wms/models"
	"wms/repositories"
	"wms/utils"

	"gorm.io/gorm"
)

type InboundBastService interface {
	ProcessBastUpload(
		supplier string,
		headerBarcode, headerName, headerItem, headerPrice string,
		fileName string,
		file multipart.File,
		fileType string,
		db *gorm.DB,
	) (inserted int, skipped int, skipDetails []string, err error)
	// ScanAndMovePendingToMaster(documentID string, db *gorm.DB) (migrated, skipped int, details []string, err error) // dihapus
	GetDocumentByID(documentID string, db *gorm.DB) (*models.ProductDocument, error)
	GetPendingProductByBarcode(documentID, barcode string, db *gorm.DB) (*models.ProductPending, error)
	ScanAndMoveSinglePendingToMaster(documentID, barcode string, categoryIDInput *string, statusInput string, note string, db *gorm.DB) (bool, string, int, error)
	GetInboundBastSummary(db *gorm.DB, date *time.Time, dateStart *time.Time, dateEnd *time.Time) (response.InboundBastSummaryResponse, error)
}

type inboundBastService struct{}

func NewInboundBastService() InboundBastService {
	return &inboundBastService{}
}

func (s *inboundBastService) GetInboundBastSummary(db *gorm.DB, date *time.Time, dateStart *time.Time, dateEnd *time.Time) (response.InboundBastSummaryResponse, error) {
	filter := repositories.InboundBastSummaryFilter{
		Date:      date,
		DateStart: dateStart,
		DateEnd:   dateEnd,
	}
	result, err := repositories.GetInboundBastSummary(db, filter)
	if err != nil {
		return response.InboundBastSummaryResponse{}, err
	}
	return response.InboundBastSummaryResponse{
		TotalFileUpload:       result.TotalFileUpload,
		TotalFileMasihProses:  result.TotalFileMasihProses,
		TotalItemTerScan:      result.TotalItemTerScan,
		TotalHargaAsalTerscan: result.TotalHargaAsalTerscan,
	}, nil
}

func (s *inboundBastService) ProcessBastUpload(
	supplier string,
	headerBarcode, headerName, headerItem, headerPrice string,
	fileName string,
	file multipart.File,
	fileType string,
	db *gorm.DB,
) (inserted int, skipped int, skipDetails []string, err error) {
	skipDetails = []string{}
	if strings.TrimSpace(fileName) == "" {
		fileName = "bast_upload"
	}

	// 1. Parse file (xlsx/csv)
	headers, rows, err := utils.ParseBulkFile(file, fileType)
	if err != nil {
		return 0, 0, []string{fmt.Sprintf("Gagal parsing file: %v", err)}, err
	}

	// Normalisasi header: trim dan lower-case agar mapping lebih fleksibel
	normalize := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}

	headerIndex := make(map[string]int)
	for idx, h := range headers {
		headerIndex[normalize(h)] = idx
	}

	idxBarcode, okBarcode := headerIndex[normalize(headerBarcode)]
	idxName, okName := headerIndex[normalize(headerName)]
	idxItem, okItem := headerIndex[normalize(headerItem)]
	idxPrice, okPrice := headerIndex[normalize(headerPrice)]
	if !okBarcode || !okName || !okItem || !okPrice {
		return 0, 0, []string{"Header mapping tidak ditemukan di file"}, fmt.Errorf("header mapping tidak ditemukan di file")
	}

	getCell := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	isEmptyDataRow := func(row []string) bool {
		return getCell(row, idxBarcode) == "" &&
			getCell(row, idxName) == "" &&
			getCell(row, idxItem) == "" &&
			getCell(row, idxPrice) == ""
	}

	fileItem := 0
	filePrice := 0
	for _, row := range rows {
		if isEmptyDataRow(row) {
			continue
		}
		fileItem++
		filePrice += int(utils.ParseFloatDefault(getCell(row, idxPrice)))
	}

	// 2. Create productDocument
	doc := models.ProductDocument{
		Code:          fmt.Sprintf("BAST-%d", utils.NowUnixNano()),
		FileName:      fileName,
		FileItem:      fileItem,
		FilePrice:     filePrice,
		Status:        "progress",
		Type:          "bast",
		HeaderBarcode: headerBarcode,
		HeaderName:    headerName,
		HeaderItem:    headerItem,
		TypeProduct:   nil,
		HeaderPrice:   headerPrice,
		Supplier:      supplier,
	}
	if err := db.Create(&doc).Error; err != nil {
		return 0, 0, []string{fmt.Sprintf("Gagal simpan dokumen: %v", err)}, err
	}

	// 3. Insert ke productPending
	for _, row := range rows {
		if isEmptyDataRow(row) {
			continue
		}
		if len(row) <= idxBarcode || len(row) <= idxName || len(row) <= idxItem || len(row) <= idxPrice {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: kolom kurang lengkap: %v", row))
			skipped++
			continue
		}
		pending := models.ProductPending{
			DocumentID: doc.ID.String(),
			Barcode:    getCell(row, idxBarcode),
			Name:       getCell(row, idxName),
			Item:       utils.ParseIntDefault(getCell(row, idxItem)),
			Price:      utils.ParseFloatDefault(getCell(row, idxPrice)),
			Status:     "discrepancy",
			IsSKU:      false,
			Note:       "",
		}
		if err := db.Create(&pending).Error; err != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: gagal insert product_pending: %v | DB error: %v", row, err))
			skipped++
			continue
		}
		inserted++
	}
	return inserted, skipped, skipDetails, nil
}

// GetDocumentByID untuk ambil data dokumen BAST
func (s *inboundBastService) GetDocumentByID(documentID string, db *gorm.DB) (*models.ProductDocument, error) {
	docRepo := repositories.NewProductDocumentRepository(db)
	docs, err := docRepo.FindAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		if doc.ID.String() == documentID {
			return &doc, nil
		}
	}
	return nil, fmt.Errorf("Dokumen tidak ditemukan")
}

// GetPendingProductByBarcode untuk ambil product pending berdasarkan barcode
func (s *inboundBastService) GetPendingProductByBarcode(documentID, barcode string, db *gorm.DB) (*models.ProductPending, error) {
	pendingRepo := repositories.NewProductPendingRepository(db)
	return pendingRepo.FindByDocumentIDAndBarcode(documentID, barcode)
}

// ScanAndMoveSinglePendingToMaster migrasi satu product pending ke master
// Sekarang menerima categoryID dan status dari user
func (s *inboundBastService) ScanAndMoveSinglePendingToMaster(documentID, barcode string, categoryIDInput *string, statusInput string, note string, db *gorm.DB) (bool, string, int, error) {
	// Cek dokumen sudah di-finish/lock
	docRepo := repositories.NewProductDocumentRepository(db)
	docs, err := docRepo.FindAll()
	if err != nil {
		return false, "Gagal mengambil dokumen", 500, err
	}
	var doc *models.ProductDocument
	for _, d := range docs {
		if d.ID.String() == documentID {
			doc = &d
			break
		}
	}
	if doc == nil {
		return false, "Dokumen tidak ditemukan", 404, nil
	}
	if doc.DateStop != nil || doc.Status == "done" {
		return false, "Dokumen sudah di-finish/lock, tidak bisa scan lagi", 409, nil
	}
	pendingRepo := repositories.NewProductPendingRepository(db)
	stickerRepo := repositories.NewStickerRepository(db)

	p, err := pendingRepo.FindByDocumentIDAndBarcode(documentID, barcode)
	if err != nil {
		return false, "Product pending tidak ditemukan", 404, err
	}
	if p.Status == "good" {
		return false, "Product sudah dipindahkan", 409, nil
	}

	var priceWarehouse float64
	var categoryID, stickerID *string

	if p.Price >= 100000 {
		// Reguler: categoryID wajib dari input user
		if categoryIDInput == nil {
			return false, "CategoryID wajib diisi untuk produk reguler", 400, nil
		}
		categoryID = categoryIDInput
		var category models.Category
		if err := db.Where("id = ?", *categoryID).First(&category).Error; err == nil && category.Discount != nil {
			discount := float64(*category.Discount)
			priceWarehouse = p.Price * (1 - discount/100)
		}
	} else if p.Price < 100000 {
		// Sticker: cari sticker berdasarkan range harga
		stickers, err := stickerRepo.List()
		if err != nil {
			return false, "Gagal mengambil data sticker", 500, err
		}
		foundSticker := false
		var sticker *models.Sticker
		for _, s := range stickers {
			if s.MinPrice != nil && s.MaxPrice != nil && p.Price >= float64(*s.MinPrice) && p.Price <= float64(*s.MaxPrice) {
				sticker = &s
				foundSticker = true
				break
			}
		}
		if !foundSticker || sticker == nil || sticker.FixedPrice == nil {
			return false, "Sticker/fixed_price tidak ditemukan untuk harga ini", 404, nil
		}
		id := sticker.ID.String()
		stickerID = &id
		priceWarehouse = float64(*sticker.FixedPrice)
	} else {
		return false, "Tidak memenuhi syarat reguler/sticker", 400, nil
	}

	var location string
	if categoryID != nil {
		location = "staging_reguler"
	} else if stickerID != nil {
		location = "staging_sticker"
	} else {
		location = ""
	}

	master := models.ProductMaster{
		DocumentID:       p.DocumentID,
		Barcode:          p.Barcode,
		Name:             p.Name,
		Item:             p.Item,
		Price:            p.Price,
		PriceWarehouse:   priceWarehouse,
		CategoryID:       categoryID,
		StickerID:        stickerID,
		ProductPendingID: func() *string { id := p.ID.String(); return &id }(),
		IsSKU:            p.IsSKU,
		Location:         location,
		TypeOut:          nil,
	}
	if err := db.Create(&master).Error; err != nil {
		return false, "Gagal insert ke master", 500, err
	}

	p.Status = statusInput
	if statusInput != "good" {
		p.Note = note
	} else {
		p.Note = ""
	}
	now := utils.Now()
	p.DateScanned = &now
	if err := pendingRepo.Update(p); err != nil {
		return true, "Berhasil migrasi, tapi gagal update status pending", 500, err
	}
	return true, "Berhasil migrasi", 200, nil
}
