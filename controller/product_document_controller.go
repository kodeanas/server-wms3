package controller

import (
	"errors"
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductDocumentController struct {
	service services.ProductDocumentService
}

func NewProductDocumentController(service services.ProductDocumentService) *ProductDocumentController {
	return &ProductDocumentController{service: service}
}

func (ctl *ProductDocumentController) ListDocuments(c *gin.Context) {
	docs, err := ctl.service.ListDocuments()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, docs, "List product documents", nil, http.StatusOK)
}

// Bulk
// GetBulkDocuments khusus untuk ambil data dengan type bulk
func (ctl *ProductDocumentController) GetBulkDocuments(c *gin.Context) {
	// Memanggil method GetBulkDocuments yang kita tambahkan di service sebelumnya
	docs, err := ctl.service.GetBulkDocuments()
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil data bulk: "+err.Error())
		return
	}

	utils.SendSuccess(c, docs, "List bulk product documents", nil, http.StatusOK)
}

func (ctl *ProductDocumentController) GetBulkDocumentDetail(c *gin.Context) {
	id := c.Param("id")
	doc, err := ctl.service.GetBulkDocumentDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(c, http.StatusNotFound, "Bulk document tidak ditemukan")
			return
		}
		utils.SendError(c, 500, "Gagal mengambil detail bulk document: "+err.Error())
		return
	}

	utils.SendSuccess(c, doc, "Detail bulk product document", nil, http.StatusOK)
}
