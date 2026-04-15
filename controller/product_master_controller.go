package controller

import (
	"errors"
	"net/http"
	dto "wms/dto/response"
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
		utils.SendError(c, 400, "Kolom rackStagingID kosong")
		return
	}
	master, err := ctl.service.GetByBarcodeWarehouse(req.BarcodeWarehouse)
	if err != nil {
		utils.SendError(c, 404, "Produk tidak ditemukan")
		return
	}
	// Filter hanya location staging_reguler
	if master.Location != "staging_reguler" {
		utils.SendError(c, 400, "Hanya produk dengan lokasi staging_reguler yang dapat di-scan ke rack staging")
		return
	}
	if master.RackStagingID != nil && *master.RackStagingID != "" {
		utils.SendError(c, 400, "Produk sudah dimasukkan ke rack staging, tidak dapat di scan ulang")
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

// List all product master in a rack staging
func (ctl *ProductMasterController) ListByRackStagingID(c *gin.Context) {
	rackStagingID := c.Param("rackStagingID")
	if rackStagingID == "" {
		utils.SendError(c, 400, "Kolom rackStagingID kosong")
		return
	}
	masters, err := ctl.service.ListByRackStagingID(rackStagingID)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, masters, "List produk di rack staging", nil, http.StatusOK)
}
