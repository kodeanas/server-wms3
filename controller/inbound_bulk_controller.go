package controller

import (
	"errors"
	"net/http"
	response "wms/dto/response"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var inboundBulkService = services.NewInboundBulkService()

// Handler untuk summary all BULK
func InboundBulkSummaryAllHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := inboundBulkService.GetBulkSummaryAll(db)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendItemSuccess(c, result, "", http.StatusOK)
	}
}

// Handler untuk detail dokumen BULK (code, nama, totalprice, total item)
func InboundBulkDocumentDetailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")

		result, err := inboundBulkService.GetBulkDocumentDetail(documentID, db)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendError(c, http.StatusNotFound, "Dokumen BULK tidak ditemukan")
				return
			}
			utils.SendError(c, 500, err.Error())
			return
		}

		utils.SendItemSuccess(c, result, "", http.StatusOK)
	}
}

// Handler untuk summary dokumen BULK (per category/sticker)
func InboundBulkDocumentSummaryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")

		result, err := inboundBulkService.GetBulkDocumentSummary(documentID, db)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendError(c, http.StatusNotFound, "Dokumen BULK tidak ditemukan")
				return
			}
			utils.SendError(c, 500, err.Error())
			return
		}

		utils.SendItemSuccess(c, result, "", http.StatusOK)
	}
}

// Handler untuk list produk berdasarkan document BULK
func InboundBulkDocumentProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")
		pg := utils.ParsePagination(c, 10)
		searchName := c.Query("name")
		searchBarcode := c.Query("barcode")

		products, total, err := inboundBulkService.GetBulkDocumentProducts(documentID, pg.Page, pg.Limit, searchName, searchBarcode, db)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendError(c, http.StatusNotFound, "Dokumen BULK tidak ditemukan")
				return
			}
			utils.SendError(c, 500, err.Error())
			return
		}

		if products == nil {
			products = make([]response.BulkProductDocumentItemResponse, 0)
		}

		utils.SendListSuccess(c, products, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}

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

		headers, rows, err := utils.ParseBulkFile(file, fileType)
		if err != nil {
			utils.SendError(c, 400, "Gagal membaca file: "+err.Error())
			return
		}

		mapping := models.BulkInboundMapping{
			BarcodeHeader: "barcode",
			NameHeader:    "name",
			QtyHeader:     "qty",
			PriceHeader:   "price",
		}
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
			FileName:    header.Filename,
			Supplier:    supplier,
			TypeProduct: typeProduct,
			Type:        fileType,
			Mapping:     mapping,
			Rows:        rows,
			Headers:     headers,
		}

		inserted, skipped, skipDetails := inboundBulkService.InboundBulkProcess(req, db)
		utils.SendSuccess(c, gin.H{
			"inserted":     inserted,
			"skipped":      skipped,
			"skip_details": skipDetails,
			"filename":     header.Filename,
		}, "Bulk inbound selesai", nil, http.StatusOK)
	}
}
