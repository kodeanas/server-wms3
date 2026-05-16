package controller

import (
	"math"
	"net/http"
	dto "wms/dto/response"
	"wms/models"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InboundDocumentSummaryHandler menampilkan summary/statistik barang berdasarkan status
func InboundDocumentSummaryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("id")

		// Ambil info dokumen
		var doc models.ProductDocument
		if err := db.Where("id = ?", documentID).First(&doc).Error; err != nil {
			utils.SendError(c, 404, "Dokumen tidak ditemukan")
			return
		}

		// Ambil semua product pending dari dokumen ini
		var pendings []models.ProductPending
		if err := db.Where("document_id = ?", documentID).Find(&pendings).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		// Hitung breakdown per status
		statusMap := make(map[string]struct {
			totalItem  int
			totalPrice float64
		})

		totalAllItem := 0
		totalAllPrice := 0.0

		for _, p := range pendings {
			status := p.Status
			if status == "" {
				status = "file"
			}

			current := statusMap[status]
			current.totalItem += 1
			current.totalPrice += p.Price
			statusMap[status] = current

			totalAllItem += 1
			totalAllPrice += p.Price
		}

		// Helper function untuk calculate persentase
		calculatePricePersentase := func(price float64) float64 {
			if totalAllPrice == 0 {
				return 0
			}
			return math.Round((price/totalAllPrice)*10000) / 100
		}

		// Build response dengan breakdown per status, tambahkan status dokumen
		response := dto.InboundDocumentSummary{
			Code:      doc.Code,
			FileName:  doc.FileName,
			FileItem:  doc.FileItem,
			FilePrice: float64(doc.FilePrice),
			Status:    doc.Status,
			Good: dto.StatusBreakdown{
				TotalItem:  statusMap["good"].totalItem,
				TotalPrice: statusMap["good"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["good"].totalPrice),
			},
			Damaged: dto.StatusBreakdown{
				TotalItem:  statusMap["damaged"].totalItem,
				TotalPrice: statusMap["damaged"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["damaged"].totalPrice),
			},
			Abnormal: dto.StatusBreakdown{
				TotalItem:  statusMap["abnormal"].totalItem,
				TotalPrice: statusMap["abnormal"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["abnormal"].totalPrice),
			},
			Non: dto.StatusBreakdown{
				TotalItem:  statusMap["non"].totalItem,
				TotalPrice: statusMap["non"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["non"].totalPrice),
			},
			Discrepancy: dto.StatusBreakdown{
				TotalItem:  statusMap["discrepancy"].totalItem,
				TotalPrice: statusMap["discrepancy"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["discrepancy"].totalPrice),
			},
			File: dto.StatusBreakdown{
				TotalItem:  statusMap["file"].totalItem,
				TotalPrice: statusMap["file"].totalPrice,
				Persentase: calculatePricePersentase(statusMap["file"].totalPrice),
			},
		}

		utils.SendItemSuccess(c, response, "", http.StatusOK)
	}
}

// InboundProductListDiscrepancyHandler menampilkan list barang yang belum terscan/status discrepancy
func InboundProductListDiscrepancyHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("id")
		pg := utils.ParsePagination(c, 10)
		searchName := c.Query("name")
		searchBarcode := c.Query("barcode")

		// Build query dengan base condition
		query := db.Model(&models.ProductPending{}).
			Where("document_id = ? AND (status = ? OR status = ? OR date_scanned IS NULL)", documentID, "discrepancy", "")

		// Tambahkan search condition (prioritas: name > barcode)
		if searchName != "" {
			query = query.Where("name LIKE ?", "%"+searchName+"%")
		} else if searchBarcode != "" {
			query = query.Where("barcode LIKE ?", "%"+searchBarcode+"%")
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		var pendings []models.ProductPending
		if err := query.Order("created_at DESC").
			Limit(pg.Limit).Offset(pg.Offset).
			Find(&pendings).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		response := make([]dto.InboundProductPendingResponse, 0)
		for _, p := range pendings {
			response = append(response, dto.InboundProductPendingResponse{
				ID:          p.ID.String(),
				Barcode:     p.Barcode,
				Name:        p.Name,
				Item:        p.Item,
				ItemGood:    p.ItemGood,
				ItemDamaged: p.ItemDamaged,
				Price:       p.Price,
				Status:      p.Status,
				Note:        p.Note,
				DateScanned: p.DateScanned,
				IsSKU:       p.IsSKU,
			})
		}

		utils.SendListSuccess(c, response, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}

// InboundProductListScannedHandler menampilkan list barang yang sudah terscan/status selain discrepancy
func InboundProductListScannedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("id")
		pg := utils.ParsePagination(c, 10)
		searchName := c.Query("name")
		searchBarcode := c.Query("barcode")

		// Build query dengan base condition
		query := db.Model(&models.ProductPending{}).
			Where("document_id = ? AND status IN (?) AND date_scanned IS NOT NULL", documentID, []string{"good", "damaged", "abnormal", "non"})

		// Tambahkan search condition (prioritas: name > barcode)
		if searchName != "" {
			query = query.Where("name LIKE ?", "%"+searchName+"%")
		} else if searchBarcode != "" {
			query = query.Where("barcode LIKE ?", "%"+searchBarcode+"%")
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		var pendings []models.ProductPending
		if err := query.Order("created_at DESC").
			Limit(pg.Limit).Offset(pg.Offset).
			Find(&pendings).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		response := make([]dto.InboundProductPendingResponse, 0)
		for _, p := range pendings {
			response = append(response, dto.InboundProductPendingResponse{
				ID:          p.ID.String(),
				Barcode:     p.Barcode,
				Name:        p.Name,
				Item:        p.Item,
				ItemGood:    p.ItemGood,
				ItemDamaged: p.ItemDamaged,
				Price:       p.Price,
				Status:      p.Status,
				Note:        p.Note,
				DateScanned: p.DateScanned,
				IsSKU:       p.IsSKU,
			})
		}

		utils.SendListSuccess(c, response, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}
