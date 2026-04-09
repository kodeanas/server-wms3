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

// Implementasi filter bast
func (ctl *ProductDocumentController) GetBastDocuments(c *gin.Context) {
	docs, err := ctl.service.GetBastDocuments()
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil data bast: "+err.Error())
		return
	}

	utils.SendSuccess(c, docs, "List bast product documents", nil, http.StatusOK)
}

func (ctl *ProductDocumentController) GetBastRelationsDetail(c *gin.Context) {
	id := c.Param("id")
	data, err := ctl.service.GetBastRelationsDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(c, http.StatusNotFound, "Bast document tidak ditemukan")
			return
		}
		utils.SendError(c, 500, "Gagal mengambil relasi bast document: "+err.Error())
		return
	}

	utils.SendSuccess(c, data, "Detail relasi bast product document", nil, http.StatusOK)
}

func (ctl *ProductDocumentController) GetBastOverview(c *gin.Context) {
	id := c.Param("id")
	overview, err := ctl.service.GetBastOverview(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(c, http.StatusNotFound, "Bast document tidak ditemukan")
			return
		}
		utils.SendError(c, 500, "Gagal mengambil overview bast document: "+err.Error())
		return
	}

	utils.SendSuccess(c, overview, "Overview bast product document", nil, http.StatusOK)
}

func (ctl *ProductDocumentController) GetBastPendingByType(c *gin.Context) {
	id := c.Param("id")
	grouped, err := ctl.service.GetBastPendingsByType(id)
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil pending bast by type: "+err.Error())
		return
	}

	utils.SendSuccess(c, grouped, "Pending bast by type", nil, http.StatusOK)
}

// Finish/lock dokumen BAST (isi date_stop)
func (ctl *ProductDocumentController) FinishDocument(c *gin.Context) {
	id := c.Param("document_id")
	err := ctl.service.FinishDocument(id)
	if err != nil {
		utils.SendError(c, 500, "Gagal finish dokumen: "+err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Dokumen berhasil di-finish/lock", nil, http.StatusOK)
}

func (c *InboundSKUController) ListSKUProductDocuments(ctx *gin.Context) {
	docs, err := c.Service.ListSKUProductDocuments()
	if err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccess(ctx, docs, "List SKU product documents", nil, http.StatusOK)
}
