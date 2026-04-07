package controller

import (
	"net/http"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var inboundBulkService = services.NewInboundBulkService()

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
