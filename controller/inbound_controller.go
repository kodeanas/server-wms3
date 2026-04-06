// ListAllProductMastersHandler menampilkan seluruh data master secara descending
package controller

import (
	"net/http"

	dto "wms/dto/response/inbound"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ListAllProductMastersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var masters []models.ProductPending
		if err := db.Order("created_at DESC").Find(&masters).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendSuccess(c, masters, "List master data", nil, http.StatusOK)
	}
}

// func ListProductManualHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var manuals []models.ProductPending
// 		if err := db.Order("created_at DESC").Find(&manuals).Error; err != nil {
// 			utils.SendError(c, 500, err.Error())
// 			return
// 		}
// 		utils.SendSuccess(c, manuals, "List manual data", nil, http.StatusOK)
// 	}
// }

func ListProductManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var results []dto.ProductManualResponse

		err := db.Table("product_pendings pp").
			Select(`
				pp.id,
				pp.document_id,
				pp.barcode,
				pp.name,
				pp.item,
				pp.price,
				pp.status,
				pp.note,
				pp.created_at,

				pm.price_warehouse,
				c.name AS category_name,
				s.name AS sticker_name
			`).
			Joins(`
				LEFT JOIN product_masters pm 
				ON pm.product_pending_id = pp.id
			`).
			Joins(`
				LEFT JOIN categories c 
				ON c.id = pm.category_id
			`).
			Joins(`
				LEFT JOIN stickers s 
				ON s.id = pm.sticker_id
			`).
			Order("pp.created_at DESC").
			Scan(&results).Error

		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		utils.SendSuccess(c, results, "List manual data", nil, http.StatusOK)
	}
}

// Tambahkan variabel global untuk service
var inboundService = services.NewInboundService()

// Handler untuk upload dan proses bulk sekaligus (single step)
func InboundBulkUploadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		supplier := c.PostForm("supplier")
		typeProduct := c.PostForm("type_product") // reguler/sticker
		fileType := c.PostForm("type")            // csv/xlsx/xls

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			utils.SendError(c, 400, "File tidak ditemukan")
			return
		}
		defer file.Close()

		// Parse file
		headers, rows, err := utils.ParseBulkFile(file, fileType)
		if err != nil {
			utils.SendError(c, 400, "Gagal membaca file: "+err.Error())
			return
		}

		// Mapping otomatis (asumsi header sudah sesuai mapping FE, atau FE kirim mapping di form-data jika perlu)
		mapping := models.BulkInboundMapping{
			BarcodeHeader: "barcode", // default, bisa diambil dari FE jika dinamis
			NameHeader:    "name",
			QtyHeader:     "qty",
			PriceHeader:   "price",
		}
		// Jika FE kirim mapping, ambil dari form-data
		if v := c.PostForm("barcode_header"); v != "" {
			mapping.BarcodeHeader = v
		}
		if v := c.PostForm("name_header"); v != "" {
			mapping.NameHeader = v
		}
		if v := c.PostForm("qty_header"); v != "" {
			mapping.QtyHeader = v
		}
		if v := c.PostForm("price_header"); v != "" {
			mapping.PriceHeader = v
		}

		req := models.BulkInboundRequest{
			Supplier:    supplier,
			TypeProduct: typeProduct,
			Type:        fileType,
			Mapping:     mapping,
			Rows:        rows,
			Headers:     headers,
		}

		inserted, skipped, skipDetails := inboundService.InboundBulkProcess(req, db)
		utils.SendSuccess(c, gin.H{
			"inserted":     inserted,
			"skipped":      skipped,
			"skip_details": skipDetails,
			"filename":     header.Filename,
		}, "Bulk inbound selesai", nil, http.StatusOK)
	}
}

func InboundManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.InboundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			var verrs []utils.ErrorItem
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, ferr := range ve {
					verrs = append(verrs, utils.ErrorItem{
						Field:   ferr.Field(),
						Message: ferr.Error(),
					})
				}
			} else {
				verrs = append(verrs, utils.ErrorItem{Field: "", Message: err.Error()})
			}
			utils.SendValidationError(c, verrs)
			return
		}

		pending, master, err := inboundService.InboundManual(models.InboundRequest{
			Name:       req.Name,
			Item:       req.Item,
			Price:      req.Price,
			CategoryID: req.CategoryID,
			StickerID:  req.StickerID,
			Status:     req.Status,
		}, db)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendSuccess(c, gin.H{"pending": pending, "master": master}, "Inbound berhasil dibuat", nil, http.StatusOK)
	}
}
