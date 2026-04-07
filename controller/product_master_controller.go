package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
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
