// ListAllProductMastersHandler menampilkan seluruh data master secara descending
package controller

import (
	"net/http"

	dto "wms/dto/response"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Tambahkan variabel global untuk service
var inboundService = services.NewInboundService(nil)

func ListAllProductMastersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg := utils.ParsePagination(c, 10)
		search := c.Query("search")

		query := db.Model(&models.ProductMaster{})
		if search != "" {
			like := "%" + search + "%"
			query = query.Where("barcode_warehouse ILIKE ? OR name_warehouse ILIKE ? OR barcode ILIKE ? OR name ILIKE ?", like, like, like, like)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		var masters []models.ProductMaster
		if err := query.Order("created_at DESC").Limit(pg.Limit).Offset(pg.Offset).Find(&masters).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendListSuccess(c, masters, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}

func ListAllProductPendingsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg := utils.ParsePagination(c, 10)
		search := c.Query("search")

		query := db.Model(&models.ProductPending{})
		if search != "" {
			like := "%" + search + "%"
			query = query.Where("barcode ILIKE ? OR name ILIKE ?", like, like)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		var pendings []models.ProductPending
		if err := query.Order("created_at DESC").Limit(pg.Limit).Offset(pg.Offset).Find(&pendings).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendListSuccess(c, pendings, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}

func ListProductManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg := utils.ParsePagination(c, 10)
		search := c.Query("search")

		base := db.Table("product_pendings pp").
			Joins("LEFT JOIN product_documents pd ON pd.id = pp.document_id").
			Joins("LEFT JOIN product_masters pm ON pm.product_pending_id = pp.id").
			Joins("LEFT JOIN categories c ON c.id = pm.category_id").
			Joins("LEFT JOIN stickers s ON s.id = pm.sticker_id").
			Where("pd.type = ?", "manual")

		if search != "" {
			like := "%" + search + "%"
			base = base.Where("pp.barcode ILIKE ? OR pp.name ILIKE ?", like, like)
		}

		var total int64
		if err := base.Count(&total).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		var results []dto.ProductManualResponse
		err := base.
			Select(`
				pp.id,
				pp.document_id,
				pp.barcode,
				pp.name,
				pp.item,
				pp.price,
				pp.status,
				pp.note,
				pp.created_at,

				pm.price_warehouse,
				c.name AS category_name,
				s.name AS sticker_name
			`).
			Order("pp.created_at DESC").
			Limit(pg.Limit).Offset(pg.Offset).
			Scan(&results).Error

		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		utils.SendListSuccess(c, results, pg.Page, pg.Limit, total, "", http.StatusOK)
	}
}

func InboundManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.InboundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			var verrs []utils.ErrorItem
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, ferr := range ve {
					verrs = append(verrs, utils.ErrorItem{
						Field:   ferr.Field(),
						Message: ferr.Error(),
					})
				}
			} else {
				verrs = append(verrs, utils.ErrorItem{Field: "", Message: err.Error()})
			}
			utils.SendValidationError(c, verrs)
			return
		}

		_, master, err := inboundService.InboundManual(models.InboundRequest{
			Name:       req.Name,
			Item:       req.Item,
			Price:      req.Price,
			CategoryID: req.CategoryID,
			StickerID:  req.StickerID,
			Status:     req.Status,
			Note:       req.Note,
		}, db)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		// Ambil nama kategori jika ada
		var categoryName interface{} = nil
		if master.CategoryID != nil {
			var category models.Category
			if err := db.Where("id = ?", *master.CategoryID).First(&category).Error; err == nil {
				categoryName = category.Name
			}
		}

		// Response hanya productMaster
		utils.SendSuccess(c, gin.H{
			"message":         "Berhasil migrasi",
			"barcode":         master.BarcodeWarehouse,
			"price":           master.Price,
			"price_warehouse": master.PriceWarehouse,
			"name":            master.Name,
			"category_name":   categoryName,
		}, "OK", nil, http.StatusOK)
	}
}
