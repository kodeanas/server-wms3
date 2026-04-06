package controller

import (
	"net/http"
	"time"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type ProductMasterSummaryController struct {
	service services.ProductMasterSummaryService
}

func NewProductMasterSummaryController(service services.ProductMasterSummaryService) *ProductMasterSummaryController {
	return &ProductMasterSummaryController{service: service}
}

func (ctl *ProductMasterSummaryController) GetSummary(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	// Perbaiki urutan parsing dan defaulting agar tidak error jika salah satu kosong
	if fromStr == "" && toStr == "" {
		utils.SendError(c, http.StatusBadRequest, "from or to date required")
		return
	}
	if fromStr == "" && toStr != "" {
		fromStr = toStr
	}
	if fromStr != "" && toStr == "" {
		toStr = fromStr
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid from date format (yyyy-mm-dd)")
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid to date format (yyyy-mm-dd)")
		return
	}
	// Set waktu to ke akhir hari agar data pada tanggal to tetap terambil
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	summary, err := ctl.service.GetSummary(from, to)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, summary, "Ringkasan pieces", nil, http.StatusOK)
}
