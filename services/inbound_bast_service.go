package services

import (
	"fmt"
	"mime/multipart"
	"strings"
	"wms/models"
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
}

type inboundBastService struct{}

func NewInboundBastService() InboundBastService {
	return &inboundBastService{}
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
