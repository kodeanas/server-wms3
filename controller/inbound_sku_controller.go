package controller

import (
	"net/http"
	importDto "wms/dto/response"
	"wms/repositories"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InboundSKUController struct {
	Service services.InboundSKUService
}

func NewInboundSKUController(service services.InboundSKUService) *InboundSKUController {
	return &InboundSKUController{Service: service}
}

func (c *InboundSKUController) UploadExcel(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "file is required")
		return
	}
	supplier := ctx.PostForm("supplier")
	if supplier == "" {
		utils.SendError(ctx, http.StatusBadRequest, "supplier is required")
		return
	}
	defer file.Close()

	// Tentukan tipe file dari ekstensi
	fileType := ""
	if len(header.Filename) > 5 && header.Filename[len(header.Filename)-5:] == ".xlsx" {
		fileType = "xlsx"
	} else if len(header.Filename) > 4 && header.Filename[len(header.Filename)-4:] == ".xls" {
		fileType = "xls"
	} else if len(header.Filename) > 4 && header.Filename[len(header.Filename)-4:] == ".csv" {
		fileType = "csv"
	} else {
		utils.SendError(ctx, http.StatusBadRequest, "unsupported file type")
		return
	}

	inserted, skipped, skipDetails, err := c.Service.UploadExcelAndCreatePendings(file, fileType, supplier, header.Filename)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(ctx, gin.H{
		"inserted":     inserted,
		"skipped":      skipped,
		"skip_details": skipDetails,
		"filename":     header.Filename,
	}, "Inbound BAST selesai", nil, http.StatusOK)
}

func InboundSKUGetDocumentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("document_id")
		docRepo := repositories.NewProductDocumentRepository(db)
		doc, err := docRepo.FindSkuDetailByID(documentID)
		if err != nil {
			utils.SendError(c, 404, "Dokumen tidak ditemukan")
			return
		}
		pendingRepo := repositories.NewProductPendingRepository(db)
		pendings, err := pendingRepo.FindByDocumentID(documentID)
		if err != nil {
			utils.SendError(c, 500, "Gagal mengambil data product pending")
			return
		}

		docDTO := importDto.InboundSKUDocumentDTO{
			ID:        doc.ID.String(),
			FileName:  doc.FileName,
			FileItem:  doc.FileItem,
			FilePrice: doc.FilePrice,
			Status:    doc.Status,
			Type:      doc.Type,
			UserID:    doc.UserID,
			Supplier:  doc.Supplier,
		}

		var productsDTO []importDto.ProductPendingDTO
		for _, p := range pendings {
			var dateScannedStr *string
			if p.DateScanned != nil {
				s := p.DateScanned.Format("2006-01-02 15:04:05")
				dateScannedStr = &s
			}
			productsDTO = append(productsDTO, importDto.ProductPendingDTO{
				ID:          p.ID.String(),
				Barcode:     p.Barcode,
				Name:        p.Name,
				Item:        p.Item,
				Price:       p.Price,
				Status:      p.Status,
				Note:        p.Note,
				DateScanned: dateScannedStr,
				ItemGood:    p.ItemGood,
				ItemDamaged: p.ItemDamaged,
			})
		}

		utils.SendSuccess(c, gin.H{
			"document":        docDTO,
			"product_pending": productsDTO,
		}, "OK", nil, 200)
	}
}

func (c *InboundSKUController) CrosscheckPending(ctx *gin.Context) {
	pendingID := ctx.Param("pending_id")
	// Cek status dokumen sebelum crosscheck
	pending, err := c.Service.GetPendingByID(pendingID)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Product pending tidak ditemukan")
		return
	}
	doc, err := c.Service.GetDocumentByID(pending.DocumentID)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Dokumen tidak ditemukan")
		return
	}
	if doc.Status == "done" {
		utils.SendError(ctx, http.StatusForbidden, "Dokumen sudah selesai dan tidak dapat di-crosscheck.")
		return
	}
	var req struct {
		ItemGood    *int `json:"item_good" binding:"required"`
		ItemDamaged *int `json:"item_damaged" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "item_good and item_damaged are required")
		return
	}
	err = c.Service.CrosscheckPending(pendingID, *req.ItemGood, *req.ItemDamaged)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(ctx, nil, "Crosscheck updated", nil, http.StatusOK)
}

func (c *InboundSKUController) FinishInboundSKU(ctx *gin.Context) {
	documentID := ctx.Param("document_id")
	err := c.Service.FinishInboundSKU(documentID)
	if err != nil {
		if err.Error() == "document already finished (done)" {
			utils.SendError(ctx, http.StatusForbidden, "Dokumen sudah selesai dan tidak dapat di-finish ulang.")
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(ctx, nil, "Inbound SKU finished", nil, http.StatusOK)
}
