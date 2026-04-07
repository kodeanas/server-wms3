package services

import (
	"fmt"
	"math/rand"
	"time"
	"wms/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InboundService interface {
	InboundManual(req models.InboundRequest, db *gorm.DB) (pending models.ProductPending, master models.ProductMaster, err error)
}

type inboundService struct{}

func NewInboundService() InboundService {
	return &inboundService{}
}

func strPtr(s string) *string {
	return &s
}

func (s *inboundService) InboundManual(req models.InboundRequest, db *gorm.DB) (models.ProductPending, models.ProductMaster, error) {
	barcode := generateUniqueBarcode()

	doc, err := getOrCreateManualDocument(db)
	if err != nil {
		return models.ProductPending{}, models.ProductMaster{}, err
	}

	var categoryID *string
	var stickerID *string
	var typeID string
	priceWarehouse := req.Price

	// =========================
	// HITUNG LOGIC
	// =========================
	if req.Price >= 100000 {
		typeID = "categories"

		if req.CategoryID != nil {
			categoryID = req.CategoryID

			var category models.Category
			if err := db.Where("id = ?", *categoryID).First(&category).Error; err == nil && category.Discount != nil {
				discount := float64(*category.Discount)
				priceWarehouse = req.Price * (1 - discount/100)
			}
		} else {
			categoryID = nil
		}

		stickerID = nil

	} else {
		typeID = "sticker"

		if req.StickerID != nil {
			stickerID = req.StickerID

			var sticker models.Sticker
			if err := db.Where("id = ?", *stickerID).First(&sticker).Error; err == nil &&
				sticker.MinPrice != nil && sticker.MaxPrice != nil {

				if req.Price >= float64(*sticker.MinPrice) &&
					req.Price <= float64(*sticker.MaxPrice) &&
					sticker.FixedPrice != nil {

					priceWarehouse = float64(*sticker.FixedPrice)
				}
			}

		} else {
			var sticker models.Sticker
			if err := db.Where("min_price <= ? AND max_price >= ?", req.Price, req.Price).
				First(&sticker).Error; err == nil && sticker.FixedPrice != nil {

				id := sticker.ID.String()
				stickerID = &id
				priceWarehouse = float64(*sticker.FixedPrice)
			} else {
				stickerID = nil
			}
		}

		categoryID = nil
	}

	// =========================
	// INSERT PENDING
	// =========================
	pending := models.ProductPending{
		ID:         uuid.New(),
		DocumentID: doc.ID.String(),
		Barcode:    barcode,
		Name:       req.Name,
		Item:       req.Item,
		Price:      req.Price,
		Status:     req.Status,
	}

	if err := db.Create(&pending).Error; err != nil {
		return pending, models.ProductMaster{}, err
	}

	// =========================
	// INSERT MASTER
	// =========================
	master := models.ProductMaster{
		DocumentID:       doc.ID.String(),
		ProductPendingID: strPtr(pending.ID.String()),
		Barcode:          barcode,
		BarcodeWarehouse: barcode,
		Name:             req.Name,
		NameWarehouse:    "Manual",
		Item:             req.Item,
		Price:            req.Price,
		PriceWarehouse:   priceWarehouse,
		CategoryID:       categoryID,
		StickerID:        stickerID,
		TypeID:           typeID,
		TypeOut:          "cargo",
	}

	// lokasi beda tergantung type
	if typeID == "categories" {
		master.Location = "staging_reguler"
	} else {
		master.Location = "staging_sticker"
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
		tp := "reguler"
		doc = models.ProductDocument{
			Code:        "INBOUND_MANUAL",
			FileName:    "INBOUND_MANUAL",
			Type:        "manual",
			Status:      "progress",
			TypeProduct: &tp,
		}
		if err := db.Create(&doc).Error; err != nil {
			return doc, err
		}
		return doc, nil
	}
	return doc, err
}
