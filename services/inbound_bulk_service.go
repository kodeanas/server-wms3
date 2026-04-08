package services

import (
	"fmt"
	"strings"
	"time"
	"wms/models"
	"wms/utils"

	"gorm.io/gorm"
)

type InboundBulkService interface {
	InboundBulkProcess(req models.BulkInboundRequest, db *gorm.DB) (inserted int, skipped int, skipDetails []string)
}

type inboundBulkService struct{}

func NewInboundBulkService() InboundBulkService {
	return &inboundBulkService{}
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
