package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var inboundBastService = services.NewInboundBastService()

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
