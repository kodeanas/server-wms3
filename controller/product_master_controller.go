package controller

import (
	"errors"
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductMasterController struct {
	service services.ProductMasterService
}

func NewProductMasterController(service services.ProductMasterService) *ProductMasterController {
	return &ProductMasterController{service: service}
}

func (ctl *ProductMasterController) ListStagingReguler(c *gin.Context) {
	masters, err := ctl.service.GetStagingReguler()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, masters, "List product master staging_reguler", nil, http.StatusOK)
}

// ListStagingSticker hanya menampilkan data dengan location = 'staging_sticker'
func (ctl *ProductMasterController) ListStagingSticker(c *gin.Context) {
	masters, err := ctl.service.GetStagingSticker()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, masters, "List product master staging_sticker", nil, http.StatusOK)
}

func (ctl *ProductMasterController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "id is required")
		return
	}

	master, err := ctl.service.GetDetailByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(c, http.StatusNotFound, "Product master not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, master, "Detail product master", nil, http.StatusOK)
}

func (ctl *ProductMasterController) UpdateStaging(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "id is required")
		return
	}

	var payload services.UpdateProductMasterStagingPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := ctl.service.UpdateStaging(id, payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(c, http.StatusNotFound, "Product master not found")
			return
		}
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, updated, "Product master staging updated", nil, http.StatusOK)
}
