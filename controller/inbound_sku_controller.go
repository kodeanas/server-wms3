package controller

import (
	"net/http"
	"wms/services"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	supplier := ctx.PostForm("supplier")
	if supplier == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "supplier is required"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
		return
	}

	err = c.Service.UploadExcelAndCreatePendings(file, fileType, supplier, header.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Upload success"})
}

func (c *InboundSKUController) CrosscheckPending(ctx *gin.Context) {
	pendingID := ctx.Param("pending_id")
	var req struct {
		ItemGood    int `json:"item_good" binding:"required"`
		ItemDamaged int `json:"item_damaged" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "item_good and item_damaged are required"})
		return
	}
	err := c.Service.CrosscheckPending(pendingID, req.ItemGood, req.ItemDamaged)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Crosscheck updated"})
}

func (c *InboundSKUController) FinishInboundSKU(ctx *gin.Context) {
	documentID := ctx.Param("document_id")
	err := c.Service.FinishInboundSKU(documentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Inbound SKU finished"})
}
