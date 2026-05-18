package controller

import (
	"net/http"
	"time"
	dto "wms/dto/response"
	"wms/models"
	"wms/repositories"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var inboundBastService = services.NewInboundBastService()

// Endpoint summary inbound BAST
func InboundBastSummaryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil query param: date, date_start, date_end
		dateStr := c.Query("date")
		dateStartStr := c.Query("date_start")
		dateEndStr := c.Query("date_end")

		var date, dateStart, dateEnd *time.Time
		layout := "2006-01-02"
		if dateStr != "" {
			d, err := time.Parse(layout, dateStr)
			if err == nil {
				date = &d
			}
		}
		if dateStartStr != "" {
			ds, err := time.Parse(layout, dateStartStr)
			if err == nil {
				dateStart = &ds
			}
		}
		if dateEndStr != "" {
			de, err := time.Parse(layout, dateEndStr)
			if err == nil {
				dateEnd = &de
			}
		}

		result, err := inboundBastService.GetInboundBastSummary(db, date, dateStart, dateEnd)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendItemSuccess(c, result, "", http.StatusOK)
	}
}

// Endpoint summary all inbound BAST
func InboundBastSummaryAllHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := inboundBastService.GetBastSummaryAll(db)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendItemSuccess(c, result, "", http.StatusOK)
	}
}

// Handler untuk upload dan proses inbound BAST
func InboundBastUploadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		supplier := c.PostForm("supplier")
		headerBarcode := c.PostForm("header_barcode")
		headerName := c.PostForm("header_name")
		headerItem := c.PostForm("header_item")
		headerPrice := c.PostForm("header_price")
		fileType := c.PostForm("type") // csv/xlsx/xls

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			utils.SendError(c, 400, "File tidak ditemukan")
			return
		}
		defer file.Close()

		inserted, skipped, skipDetails, err := inboundBastService.ProcessBastUpload(
			supplier,
			headerBarcode, headerName, headerItem, headerPrice,
			header.Filename,
			file, fileType, db,
		)
		if err != nil {
			utils.SendError(c, 400, err.Error())
			return
		}
		utils.SendSuccess(c, gin.H{
			"inserted":     inserted,
			"skipped":      skipped,
			"skip_details": skipDetails,
			"filename":     header.Filename,
		}, "Inbound BAST selesai", nil, http.StatusOK)
	}
}

// 1. Endpoint untuk masuk ke halaman scanner (ambil dokumen by id)
func InboundBastGetDocumentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")
		doc, err := inboundBastService.GetDocumentByID(documentID, db)
		if err != nil {
			utils.SendError(c, 404, err.Error())
			return
		}
		// Hitung jumlah scanned dan unscanned
		pendingRepo := repositories.NewProductPendingRepository(db)
		pendings, err := pendingRepo.FindByDocumentID(documentID)
		if err != nil {
			utils.SendError(c, 500, "Gagal mengambil data product pending")
			return
		}
		scanned := 0
		unscanned := 0
		for _, p := range pendings {
			if p.DateScanned != nil {
				scanned++
			} else {
				unscanned++
			}
		}
		// Custom document response (optional: select only needed fields)
		// Map document to DTO
		var docResp *dto.BastDocumentResponse
		if doc != nil {
			docResp = &dto.BastDocumentResponse{
				ID:        doc.ID.String(),
				Code:      doc.Code,
				FileName:  doc.FileName,
				FileItem:  doc.FileItem,
				FilePrice: doc.FilePrice,
				Status:    doc.Status,
				UserID:    doc.UserID,
				CreatedAt: doc.CreatedAt.Format(time.RFC3339Nano),
				UpdatedAt: doc.UpdatedAt.Format(time.RFC3339Nano),
				DeletedAt: doc.DeletedAt,
				DateStop:  doc.DateStop,
			}
		}
		// Map products to DTO
		var products []dto.BastProductResponse
		for _, p := range pendings {
			products = append(products, dto.BastProductResponse{
				ID:          p.ID.String(),
				Barcode:     p.Barcode,
				Name:        p.Name,
				Item:        p.Item,
				Price:       p.Price,
				Note:        p.Note,
				DateScanned: p.DateScanned,
			})
		}
		utils.SendItemSuccess(c, gin.H{
			"document":        docResp,
			"scanned_count":   scanned,
			"unscanned_count": unscanned,
			"products":        products,
		}, "", http.StatusOK)
	}
}

// 2. Endpoint untuk ambil product pending by barcode (scanner)
func InboundBastGetPendingProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")
		barcode := c.Param("barcode")
		product, err := inboundBastService.GetPendingProductByBarcode(documentID, barcode, db)
		if err != nil {
			utils.SendError(c, 404, err.Error())
			return
		}
		if product.Status == "good" {
			utils.SendError(c, 409, "Produk sudah di-scan dan tidak dapat diakses lagi.")
			return
		}
		utils.SendItemSuccess(c, product, "", http.StatusOK)
	}
}

// 3. Endpoint untuk post hasil scan (migrasi 1 produk)
func InboundBastScanSingleProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")
		barcode := c.Param("barcode")

		var req struct {
			CategoryID *string `json:"category_id"` // wajib untuk reguler
			Status     string  `json:"status"`      // good, abnormal, damaged, non
			Note       string  `json:"note"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.SendError(c, 400, "Invalid request body")
			return
		}

		success, msg, barcodeWarehouse, code, err := inboundBastService.ScanAndMoveSinglePendingToMaster(documentID, barcode, req.CategoryID, req.Status, req.Note, db)
		if err != nil || !success {
			if code == 0 {
				code = 400
			}
			utils.SendError(c, code, msg+func() string {
				if err != nil {
					return ": " + err.Error()
				} else {
					return ""
				}
			}())
			return
		}

		// Ambil data master terbaru setelah migrasi
		var master models.ProductMaster
		if err := db.Where("barcode_warehouse = ? AND document_id = ?", barcodeWarehouse, documentID).First(&master).Error; err != nil {
			utils.SendError(c, 500, "Gagal mengambil data master")
			return
		}

		// Ambil nama kategori jika ada
		var categoryName interface{} = nil
		if master.CategoryID != nil {
			var category models.Category
			if err := db.Where("id = ?", *master.CategoryID).First(&category).Error; err == nil {
				categoryName = category.Name
			}
		}

		utils.SendSuccess(c, gin.H{
			"message":         "Berhasil migrasi",
			"barcode":         master.BarcodeWarehouse,
			"price":           master.Price,
			"price_warehouse": master.PriceWarehouse,
			"name":            master.Name,
			"category_name":   categoryName,
		}, "OK", nil, http.StatusOK)
	}
}
