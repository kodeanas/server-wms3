package controller

import (
	"errors"
	"net/http"
	"time"
	dto "wms/dto/response"
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
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	docs, total, err := ctl.service.ListDocumentsPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendListSuccess(c, docs, pg.Page, pg.Limit, total, "", http.StatusOK)
}

// Bulk
// GetBulkDocuments khusus untuk ambil data dengan type bulk
func (ctl *ProductDocumentController) GetBulkDocuments(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	docs, total, err := ctl.service.GetBulkDocumentsPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil data bulk: "+err.Error())
		return
	}

	utils.SendListSuccess(c, docs, pg.Page, pg.Limit, total, "", http.StatusOK)
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

	utils.SendItemSuccess(c, doc, "", http.StatusOK)
}

// Implementasi filter bast
func (ctl *ProductDocumentController) GetBastDocuments(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	docs, total, err := ctl.service.GetBastDocumentsPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil data bast: "+err.Error())
		return
	}

	// Custom response: only selected fields for each document
	resp := make([]dto.BastDocumentResponse, 0, len(docs))
	for _, doc := range docs {
		resp = append(resp, dto.BastDocumentResponse{
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
		})
	}
	utils.SendListSuccess(c, resp, pg.Page, pg.Limit, total, "", http.StatusOK)
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

	utils.SendItemSuccess(c, data, "", http.StatusOK)
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

	utils.SendItemSuccess(c, overview, "", http.StatusOK)
}

func (ctl *ProductDocumentController) GetBastPendingByType(c *gin.Context) {
	id := c.Param("id")
	grouped, err := ctl.service.GetBastPendingsByType(id)
	if err != nil {
		utils.SendError(c, 500, "Gagal mengambil pending bast by type: "+err.Error())
		return
	}

	utils.SendItemSuccess(c, grouped, "", http.StatusOK)
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
	pg := utils.ParsePagination(ctx, 10)
	search := ctx.Query("search")

	docs, total, err := c.Service.ListSKUProductDocumentsPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendListSuccess(ctx, docs, pg.Page, pg.Limit, total, "", http.StatusOK)
}
