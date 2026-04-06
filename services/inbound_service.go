package services

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"wms/models"
	"wms/utils"

	"gorm.io/gorm"
)

type InboundService interface {
	InboundManual(req models.InboundRequest, db *gorm.DB) (pending models.ProductPending, master models.ProductMaster, err error)
	// Bulk
	InboundBulkProcess(req models.BulkInboundRequest, db *gorm.DB) (inserted int, skipped int, skipDetails []string)
}

type inboundService struct{}

func NewInboundService() InboundService {
	return &inboundService{}
}

// Proses bulk upload
func (s *inboundService) InboundBulkProcess(req models.BulkInboundRequest, db *gorm.DB) (inserted int, skipped int, skipDetails []string) {
	skipDetails = []string{}
	// 1. Simpan dokumen ke tabel ProductDocument (tetap satu kali per upload)
	doc := models.ProductDocument{
		Code:          fmt.Sprintf("BULK-%d", time.Now().UnixNano()),
		FileName:      "bulk_upload",
		Status:        "progress",
		Type:          "bulk",
		HeaderBarcode: req.Mapping.BarcodeHeader,
		HeaderName:    req.Mapping.NameHeader,
		HeaderItem:    req.Mapping.QtyHeader,
		HeaderPrice:   req.Mapping.PriceHeader,
		Supplier:      req.Supplier,
		TypeProduct:   req.TypeProduct,
		UserID:        nil, // bisa diambil dari context jika ada auth
	}
	if err := db.Create(&doc).Error; err != nil {
		return 0, 0, []string{fmt.Sprintf("Gagal simpan dokumen: %v", err)}
	}

	// 2. Hardcode mapping index kolom sesuai urutan excel
	var idxBarcode, idxName, idxCategory, idxQty, idxPrice int
	if req.TypeProduct == "sticker" {
		// Barcode, deskripsi, qty, unit price
		idxBarcode = 0
		idxName = 1
		idxQty = 2
		idxPrice = 3
		idxCategory = -1 // tidak ada kolom kategori
	} else {
		// barcode, description, category, qty, unit price, bast, discount, price after discount
		idxBarcode = 0
		idxName = 1
		idxCategory = 2
		idxQty = 3
		idxPrice = 4
	}
	// Tidak perlu validasi mapping, asumsikan urutan selalu benar

	// Ambil semua kategori dan sticker dari DB
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

	// 3. Proses setiap row: paksa tipe produk sesuai pilihan user, validasi, insert jika valid
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
		price, err2 := utils.ParseFloat(priceStr)
		if err1 != nil || err2 != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: qty/price tidak valid: %v", row))
			skipped++
			continue
		}

		var categoryID, stickerID, location, typeID string
		categoryID = ""
		stickerID = ""
		location = ""
		typeID = ""

		if req.TypeProduct == "reguler" {
			// Validasi harga
			if price < 100000 {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: harga kurang dari 100rb untuk reguler: %v", row))
				skipped++
				continue
			}
			// Validasi kategori
			kategoriNama := ""
			if idxCategory != -1 && len(row) > idxCategory {
				kategoriNama = strings.TrimSpace(row[idxCategory])
			}
			foundCategory := false
			for _, cat := range categories {
				if strings.EqualFold(strings.TrimSpace(cat.Name), kategoriNama) {
					foundCategory = true
					categoryID = cat.ID.String()
					break
				}
			}
			if !foundCategory || categoryID == "" {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: kategori tidak ditemukan di DB: '%s' (row: %v)", kategoriNama, row))
				skipped++
				continue
			}
			location = "staging_reguler"
			typeID = "categories"
		} else if req.TypeProduct == "sticker" {
			// Validasi harga
			if price >= 100000 {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: harga >= 100rb untuk sticker: %v", row))
				skipped++
				continue
			}
			// Validasi sticker range
			foundSticker := false
			for _, sticker := range stickers {
				if sticker.MinPrice != nil && sticker.MaxPrice != nil && price >= float64(*sticker.MinPrice) && price <= float64(*sticker.MaxPrice) {
					stickerID = sticker.ID.String()
					foundSticker = true
					break
				}
			}
			if !foundSticker || stickerID == "" {
				skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: tidak ada sticker dengan range harga sesuai: %v", row))
				skipped++
				continue
			}
			location = "staging_sticker"
			typeID = "sticker"
		}

		var categoryIDPtr, stickerIDPtr *string
		if categoryID == "" {
			categoryIDPtr = nil
		} else {
			categoryIDPtr = &categoryID
		}
		if stickerID == "" {
			stickerIDPtr = nil
		} else {
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
		var productPendingIDPtr *string
		if productPendingID == "" {
			productPendingIDPtr = nil
		} else {
			productPendingIDPtr = &productPendingID
		}
		master := models.ProductMaster{
			DocumentID:       doc.ID.String(),
			Barcode:          barcode,
			BarcodeWarehouse: barcode,
			Name:             name,
			NameWarehouse:    name,
			Item:             qty,
			ItemWarehouse:    0,
			Price:            price,
			PriceWarehouse:   0,
			CategoryID:       categoryIDPtr,
			StickerID:        stickerIDPtr,
			ProductPendingID: productPendingIDPtr,
			Location:         location,
			TypeID:           typeID,
			TypeOut:          "cargo",
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

func (s *inboundService) InboundManual(req models.InboundRequest, db *gorm.DB) (models.ProductPending, models.ProductMaster, error) {
	// Logic generate barcode
	barcode := generateUniqueBarcode()
	barcodeWarehouse := generateUniqueBarcode()

	// Logic dokumen: cari/buat dokumen manual
	doc, err := getOrCreateManualDocument(db)
	if err != nil {
		return models.ProductPending{}, models.ProductMaster{}, err
	}

	// Logic BE: tentukan category_id/sticker_id otomatis dan PriceWarehouse
	var categoryID, stickerID, typeID string
	var priceWarehouse float64 = req.Price

	// Insert ke ProductPending terlebih dahulu untuk dapatkan ID
	pending := models.ProductPending{
		DocumentID: doc.ID.String(),
		Barcode:    barcode,
		Name:       req.Name,
		Item:       req.Item,
		Price:      req.Price,
		Status:     req.Status, // default status valid
	}
	if err := db.Create(&pending).Error; err != nil {
		return pending, models.ProductMaster{}, err
	}

	var productPendingIDPtr *string
	productPendingID := pending.ID.String()
	if productPendingID == "" {
		productPendingIDPtr = nil
	} else {
		productPendingIDPtr = &productPendingID
	}

	var categoryIDPtr, stickerIDPtr *string

	var master models.ProductMaster
	if req.Price >= 100000 {
		if req.CategoryID != nil && *req.CategoryID != "" {
			categoryID = *req.CategoryID
			categoryIDPtr = &categoryID
			// Ambil diskon kategori
			var category models.Category
			if err := db.Where("id = ?", categoryID).First(&category).Error; err == nil && category.Discount != nil {
				discount := float64(*category.Discount)
				priceWarehouse = req.Price * (1 - discount/100)
			}
		} else {
			categoryIDPtr = nil
		}
		stickerIDPtr = nil
		typeID = "categories"

		master = models.ProductMaster{
			DocumentID:       doc.ID.String(),
			Barcode:          barcode,
			BarcodeWarehouse: barcodeWarehouse,
			Name:             req.Name,
			NameWarehouse:    "Manual",
			Item:             req.Item,
			Price:            req.Price,
			PriceWarehouse:   priceWarehouse,
			CategoryID:       categoryIDPtr,
			StickerID:        stickerIDPtr,
			ProductPendingID: productPendingIDPtr,
			TypeID:           typeID,
			Location:         "staging_reguler",
			TypeOut:          "cargo",
		}
	} else {
		if req.StickerID != nil && *req.StickerID != "" {
			stickerID = *req.StickerID
			stickerIDPtr = &stickerID
			// Cari sticker sesuai range harga
			var sticker models.Sticker
			if err := db.Where("id = ?", stickerID).First(&sticker).Error; err == nil && sticker.MinPrice != nil && sticker.MaxPrice != nil {
				if req.Price >= float64(*sticker.MinPrice) && req.Price <= float64(*sticker.MaxPrice) && sticker.FixedPrice != nil {
					priceWarehouse = float64(*sticker.FixedPrice)
				}
			}
		} else {
			// Jika stickerID tidak ada, cari sticker yang cocok dengan range harga
			var sticker models.Sticker
			if err := db.Where("min_price <= ? AND max_price >= ?", req.Price, req.Price).First(&sticker).Error; err == nil && sticker.FixedPrice != nil {
				stickerID = sticker.ID.String()
				stickerIDPtr = &stickerID
				priceWarehouse = float64(*sticker.FixedPrice)
			} else {
				stickerIDPtr = nil
			}
		}
		categoryIDPtr = nil
		typeID = "sticker"

		master = models.ProductMaster{
			DocumentID:       doc.ID.String(),
			Barcode:          barcode,
			BarcodeWarehouse: barcodeWarehouse,
			Name:             req.Name,
			NameWarehouse:    "Manual",
			Item:             req.Item,
			Price:            req.Price,
			PriceWarehouse:   priceWarehouse,
			CategoryID:       categoryIDPtr,
			StickerID:        stickerIDPtr,
			ProductPendingID: productPendingIDPtr,
			TypeID:           typeID,
			Location:         "staging_sticker",
			TypeOut:          "cargo",
		}
	}

	if err := db.Create(&master).Error; err != nil {
		return pending, master, err
	}

	return pending, master, nil
}

// Helper: generate barcode dan dokumen, bisa diambil dari controller lama
func generateUniqueBarcode() string {
	t := time.Now().UnixNano()
	r := rand.Intn(100000)
	return fmt.Sprintf("BC-%d-%d", t, r)
}

func getOrCreateManualDocument(db *gorm.DB) (models.ProductDocument, error) {
	var doc models.ProductDocument
	err := db.Where("code = ?", "INBOUND_MANUAL").First(&doc).Error
	if err == gorm.ErrRecordNotFound {
		doc = models.ProductDocument{
			Code:        "INBOUND_MANUAL",
			FileName:    "INBOUND_MANUAL",
			Type:        "manual",
			Status:      "progress",
			TypeProduct: "reguler",
		}
		if err := db.Create(&doc).Error; err != nil {
			return doc, err
		}
		return doc, nil
	}
	return doc, err
}
