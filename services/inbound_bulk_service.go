package services

import (
	"fmt"
	"strings"
	"time"
	response "wms/dto/response"
	"wms/models"
	"wms/repositories"
	"wms/utils"

	"gorm.io/gorm"
)

type InboundBulkService interface {
	InboundBulkProcess(req models.BulkInboundRequest, db *gorm.DB) (inserted int, skipped int, skipDetails []string)
	GetBulkSummaryAll(db *gorm.DB) (response.BulkSummaryAllResponse, error)
	GetBulkDocumentDetail(documentID string, db *gorm.DB) (response.BulkDocumentDetailResponse, error)
	GetBulkDocumentSummary(documentID string, db *gorm.DB) ([]response.BulkDocumentSummaryItemResponse, error)
	GetBulkDocumentProducts(documentID string, page, limit int, searchName, searchBarcode string, db *gorm.DB) ([]response.BulkProductDocumentItemResponse, int64, error)
}

type inboundBulkService struct{}

func NewInboundBulkService() InboundBulkService {
	return &inboundBulkService{}
}

func (s *inboundBulkService) GetBulkSummaryAll(db *gorm.DB) (response.BulkSummaryAllResponse, error) {
	repo := repositories.NewProductDocumentRepository(db)
	summaryMap, err := repo.GetBulkSummaryAll()
	if err != nil {
		return response.BulkSummaryAllResponse{}, err
	}
	return response.BulkSummaryAllResponse{
		TotalDocumentUpload: int(summaryMap["total_document_upload"].(int64)),
		TotalProductMasuk:   int(summaryMap["total_product_masuk"].(int64)),
		TotalHargaMasuk:     summaryMap["total_harga_masuk"].(float64),
	}, nil
}

func (s *inboundBulkService) GetBulkDocumentDetail(documentID string, db *gorm.DB) (response.BulkDocumentDetailResponse, error) {
	repo := repositories.NewProductDocumentRepository(db)
	doc, err := repo.FindBulkDetailByID(documentID)
	if err != nil {
		return response.BulkDocumentDetailResponse{}, err
	}

	return response.BulkDocumentDetailResponse{
		ID:         doc.ID.String(),
		Code:       doc.Code,
		Nama:       doc.FileName,
		TotalPrice: float64(doc.FilePrice),
		TotalItem:  doc.FileItem,
	}, nil
}

func (s *inboundBulkService) GetBulkDocumentSummary(documentID string, db *gorm.DB) ([]response.BulkDocumentSummaryItemResponse, error) {
	repo := repositories.NewProductDocumentRepository(db)
	if _, err := repo.FindBulkDetailByID(documentID); err != nil {
		return nil, err
	}
	return s.getBulkDocumentSummaryRows(documentID, db)
}

func (s *inboundBulkService) getBulkDocumentSummaryRows(documentID string, db *gorm.DB) ([]response.BulkDocumentSummaryItemResponse, error) {
	type bulkSummaryRow struct {
		Label          string  `gorm:"column:label"`
		Item           int64   `gorm:"column:item"`
		Price          float64 `gorm:"column:price"`
		PriceWarehouse float64 `gorm:"column:price_warehouse"`
	}

	rows := make([]bulkSummaryRow, 0)
	err := db.Table("product_masters pm").
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
		Where("pm.document_id = ? AND pm.deleted_at IS NULL", documentID).
		Group("label").
		Order("label ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]response.BulkDocumentSummaryItemResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, response.BulkDocumentSummaryItemResponse{
			Label:          row.Label,
			Item:           int(row.Item),
			Price:          row.Price,
			PriceWarehouse: row.PriceWarehouse,
		})
	}

	return result, nil
}

func (s *inboundBulkService) GetBulkDocumentProducts(documentID string, page, limit int, searchName, searchBarcode string, db *gorm.DB) ([]response.BulkProductDocumentItemResponse, int64, error) {
	repo := repositories.NewProductDocumentRepository(db)
	if _, err := repo.FindBulkDetailByID(documentID); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := db.Model(&models.ProductPending{}).
		Where("document_id = ?", documentID)

	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	} else if searchBarcode != "" {
		query = query.Where("barcode LIKE ?", "%"+searchBarcode+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var pendings []models.ProductPending
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&pendings).Error; err != nil {
		return nil, 0, err
	}

	result := make([]response.BulkProductDocumentItemResponse, 0, len(pendings))
	for _, p := range pendings {
		result = append(result, response.BulkProductDocumentItemResponse{
			ID:          p.ID.String(),
			Barcode:     p.Barcode,
			Name:        p.Name,
			Item:        p.Item,
			Price:       p.Price,
			Status:      p.Status,
			Note:        p.Note,
			DateScanned: p.DateScanned,
		})
	}

	return result, total, nil
}

func (s *inboundBulkService) InboundBulkProcess(req models.BulkInboundRequest, db *gorm.DB) (inserted int, skipped int, skipDetails []string) {
	skipDetails = []string{}

	var idxBarcode, idxName, idxCategory, idxQty, idxPrice int
	if req.TypeProduct == "sticker" {
		idxBarcode = 0
		idxName = 1
		idxQty = 2
		idxPrice = 3
		idxCategory = -1
	} else {
		idxBarcode = 0
		idxName = 1
		idxCategory = 2
		idxQty = 3
		idxPrice = 4
	}

	var totalFileItem int
	var totalFilePrice float64
	for _, row := range req.Rows {
		if len(row) > idxQty && len(row) > idxPrice {
			qty, _ := utils.ParseInt(row[idxQty])
			price, _ := utils.ParseCurrency(row[idxPrice])

			totalFileItem += qty
			totalFilePrice += price
		}
	}

	fileName := strings.TrimSpace(req.FileName)
	if fileName == "" {
		fileName = "bulk_upload"
	}

	typeProduct := strings.TrimSpace(req.TypeProduct)
	doc := models.ProductDocument{
		Code:          fmt.Sprintf("BULK-%d", time.Now().UnixNano()),
		FileName:      fileName,
		FileItem:      totalFileItem,
		FilePrice:     int(totalFilePrice),
		Status:        "progress",
		Type:          "bulk",
		HeaderBarcode: req.Mapping.BarcodeHeader,
		HeaderName:    req.Mapping.NameHeader,
		HeaderItem:    req.Mapping.QtyHeader,
		HeaderPrice:   req.Mapping.PriceHeader,
		Supplier:      req.Supplier,
		TypeProduct:   &typeProduct,
		UserID:        nil,
	}

	if err := db.Create(&doc).Error; err != nil {
		return 0, 0, []string{fmt.Sprintf("Gagal simpan dokumen: %v", err)}
	}

	var categories []models.Category
	var stickers []models.Sticker
	if req.TypeProduct == "reguler" {
		if err := db.Find(&categories).Error; err != nil {
			return 0, 0, []string{fmt.Sprintf("Gagal mengambil kategori: %v", err)}
		}
	} else if req.TypeProduct == "sticker" {
		if err := db.Find(&stickers).Error; err != nil {
			return 0, 0, []string{fmt.Sprintf("Gagal mengambil sticker: %v", err)}
		}
	}

	for _, row := range req.Rows {
		if len(row) <= idxPrice || len(row) <= idxQty || len(row) <= idxName || len(row) <= idxBarcode {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: kolom kurang lengkap: %v", row))
			skipped++
			continue
		}

		barcode := row[idxBarcode]
		name := row[idxName]
		qtyStr := row[idxQty]
		priceStr := row[idxPrice]

		qty, err1 := utils.ParseInt(qtyStr)
		price, err2 := utils.ParseCurrency(priceStr)
		if err1 != nil || err2 != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: qty/price tidak valid: %v", row))
			skipped++
			continue
		}

		var categoryID, stickerID, location, typeID string
		priceWarehouse := price

		if req.TypeProduct == "reguler" {
			if price < 100000 {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: type_product reguler tapi harga di bawah 100rb: %v", row))
				skipped++
				continue
			}

			kategoriNama := ""
			if idxCategory != -1 && len(row) > idxCategory {
				kategoriNama = strings.TrimSpace(row[idxCategory])
			}

			discount := 0.0
			foundCategory := false
			for _, cat := range categories {
				if strings.EqualFold(strings.TrimSpace(cat.Name), kategoriNama) {
					categoryID = cat.ID.String()
					if cat.Discount != nil {
						discount = float64(*cat.Discount)
					}
					foundCategory = true
					break
				}
			}
			if !foundCategory || categoryID == "" {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: kategori tidak ditemukan di DB: '%s' (row: %v)", kategoriNama, row))
				skipped++
				continue
			}
			priceWarehouse = price * (1 - discount/100)

			location = "staging_reguler"
			typeID = "categories"
		} else if req.TypeProduct == "sticker" {
			if price >= 100000 {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: type_product sticker tapi harga di atas/sama dengan 100rb: %v", row))
				skipped++
				continue
			}

			fixedPrice := 0.0
			foundSticker := false
			for _, sticker := range stickers {
				if sticker.MinPrice != nil && sticker.MaxPrice != nil && price >= float64(*sticker.MinPrice) && price <= float64(*sticker.MaxPrice) {
					stickerID = sticker.ID.String()
					if sticker.FixedPrice != nil {
						fixedPrice = float64(*sticker.FixedPrice)
						foundSticker = true
					}
					break
				}
			}
			if !foundSticker || stickerID == "" {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: tidak ada sticker dengan range harga sesuai: %v", row))
				skipped++
				continue
			}
			priceWarehouse = fixedPrice

			location = "staging_sticker"
			typeID = "sticker"
		}

		var categoryIDPtr, stickerIDPtr *string
		if categoryID != "" {
			categoryIDPtr = &categoryID
		}
		if stickerID != "" {
			stickerIDPtr = &stickerID
		}

		pending := models.ProductPending{
			DocumentID: doc.ID.String(),
			Barcode:    barcode,
			Name:       name,
			Item:       qty,
			Price:      price,
			Status:     "good",
			IsSKU:      false,
			Note:       "",
		}
		if err := db.Create(&pending).Error; err != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: gagal insert product_pending: %v", row))
			skipped++
			continue
		}

		productPendingID := pending.ID.String()
		master := models.ProductMaster{
			DocumentID:       doc.ID.String(),
			Barcode:          barcode,
			BarcodeWarehouse: barcode,
			Name:             name,
			NameWarehouse:    name,
			Item:             qty,
			ItemWarehouse:    qty,
			Price:            price,
			PriceWarehouse:   priceWarehouse,
			CategoryID:       categoryIDPtr,
			StickerID:        stickerIDPtr,
			ProductPendingID: &productPendingID,
			Location:         location,
			TypeID:           typeID,
			TypeOut:          nil,
		}
		if err := db.Create(&master).Error; err != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: gagal insert product_master: %v | DB error: %v", row, err))
			skipped++
			continue
		}
		inserted++
	}

	return inserted, skipped, skipDetails
}
