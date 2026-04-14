package controller

import (
	"net/http"
	dto "wms/dto/response"
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

// Scan product master by barcode_warehouse and assign to rack staging
func (ctl *ProductMasterController) ScanBarcodeWarehouse(c *gin.Context) {
	rackStagingID := c.Param("rackStagingID")
	var req struct {
		BarcodeWarehouse string `json:"barcode_warehouse" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	if rackStagingID == "" {
		utils.SendError(c, 400, "Missing rackStagingID parameter")
		return
	}
	master, err := ctl.service.GetByBarcodeWarehouse(req.BarcodeWarehouse)
	if err != nil {
		utils.SendError(c, 404, "Product not found")
		return
	}
	// Assign product to rack staging
	err = ctl.service.SetRackStaging(master.ID.String(), rackStagingID)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	// Refresh master
	master, _ = ctl.service.GetByBarcodeWarehouse(req.BarcodeWarehouse)
	resp := dto.ProductMasterScanResponse{
		ID:               master.ID.String(),
		BarcodeWarehouse: master.BarcodeWarehouse,
		NameWarehouse:    master.NameWarehouse,
		ItemWarehouse:    master.ItemWarehouse,
		Location:         master.Location,
		RackStagingID:    master.RackStagingID,
	}
	utils.SendSuccess(c, resp, "Product assigned to rack staging", nil, http.StatusOK)
}
