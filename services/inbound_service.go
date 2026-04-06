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
		UserID:        "USER ID", // bisa diambil dari context jika ada auth
	}
	if err := db.Create(&doc).Error; err != nil {
		return 0, 0, []string{fmt.Sprintf("Gagal simpan dokumen: %v", err)}
	}

	// 2. Mapping header index
	idxBarcode, idxName, idxQty, idxPrice, idxCategory := -1, -1, -1, -1, -1
	for i, h := range req.Headers {
		if h == req.Mapping.BarcodeHeader {
			idxBarcode = i
		}
		if h == req.Mapping.NameHeader {
			idxName = i
		}
		if h == req.Mapping.QtyHeader {
			idxQty = i
		}
		if h == req.Mapping.PriceHeader {
			idxPrice = i
		}
		if strings.ToLower(h) == "category" || strings.ToLower(h) == "kategori" {
			idxCategory = i
		}
	}
	if idxBarcode == -1 || idxName == -1 || idxQty == -1 || idxPrice == -1 {
		return 0, 0, []string{fmt.Sprintf("Header mapping tidak valid")}
	}

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

		master := models.ProductMaster{
			DocumentID:       doc.ID.String(),
			Barcode:          barcode,
			BarcodeWarehouse: barcode,
			Name:             name,
			NameWarehouse:    "",
			Item:             qty,
			ItemWarehouse:    0,
			Price:            price,
			PriceWarehouse:   0,
			CategoryID:       categoryID,
			StickerID:        stickerID,
			Location:         location,
			TypeID:           typeID,
			TypeOut:          "cargo",
		}
		pending := models.ProductPending{
			DocumentID: doc.ID.String(),
			Barcode:    barcode,
			Name:       name,
			Item:       qty,
			Price:      price,
			Status:     "GOOD",
			IsSKU:      false,
			Note:       "",
		}
		if err := db.Create(&master).Error; err != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: gagal insert product_master: %v | DB error: %v", row, err))
			skipped++
			continue
		}
		if err := db.Create(&pending).Error; err != nil {
			skipDetails = append(skipDetails, fmt.Sprintf("Row skipped: gagal insert product_pending: %v", row))
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
	var master models.ProductMaster
	if req.Price >= 100000 {
		if req.CategoryID != nil {
			categoryID = *req.CategoryID
			// Ambil diskon kategori
			var category models.Category
			if err := db.Where("id = ?", categoryID).First(&category).Error; err == nil && category.Discount != nil {
				discount := float64(*category.Discount)
				priceWarehouse = req.Price * (1 - discount/100)
			}
		}
		stickerID = ""
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
			CategoryID:       categoryID,
			StickerID:        stickerID,
			TypeID:           typeID,
			Location:         "staging_reguler",
			TypeOut:          "cargo",
		}
	} else {
		if req.StickerID != nil {
			stickerID = *req.StickerID
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
				priceWarehouse = float64(*sticker.FixedPrice)
			}
		}
		categoryID = ""
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
			CategoryID:       categoryID,
			StickerID:        stickerID,
			TypeID:           typeID,
			Location:         "staging_sticker",
			TypeOut:          "cargo",
		}
	}

	// Insert ke ProductPending
	pending := models.ProductPending{
		DocumentID: doc.ID.String(),
		Barcode:    barcode,
		Name:       req.Name,
		Item:       req.Item,
		Price:      req.Price,
		Status:     req.Status, // default status valid
	}
	if err := db.Create(&pending).Error; err != nil {
		return pending, master, err
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
