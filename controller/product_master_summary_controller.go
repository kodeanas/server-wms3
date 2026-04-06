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
	var from, to time.Time
	var err error
	if fromStr == "" && toStr == "" {
		// Jika keduanya kosong, pakai hari ini
		now := time.Now()
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to = from.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		if fromStr == "" && toStr != "" {
			fromStr = toStr
		}
		if fromStr != "" && toStr == "" {
			toStr = fromStr
		}
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			utils.SendError(c, http.StatusBadRequest, "Invalid from date format (yyyy-mm-dd)")
			return
		}
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			utils.SendError(c, http.StatusBadRequest, "Invalid to date format (yyyy-mm-dd)")
			return
		}
		to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}
	summary, err := ctl.service.GetSummary(from, to)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	// Jika summary kosong, balikan data null/{}
	if summary.TotalPieces == 0 && summary.TotalHargaAsal == 0 && summary.TotalHargaGudang == 0 {
		utils.SendSuccess(c, nil, "Ringkasan Pieces", nil, http.StatusOK)
		return
	}
	utils.SendSuccess(c, summary, "Ringkasan pieces", nil, http.StatusOK)
}
