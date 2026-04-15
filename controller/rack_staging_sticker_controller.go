package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type RackStagingStickerController struct {
	service services.RackStagingStickerService
}

func NewRackStagingStickerController(service services.RackStagingStickerService) *RackStagingStickerController {
	return &RackStagingStickerController{service: service}
}

// Create rackStagingSticker (auto generate code, type=sticker, is_moved=false)
func (ctl *RackStagingStickerController) Create(c *gin.Context) {
	// Ambil userID dari context (atau session, sesuaikan kebutuhan)
	userID := c.GetString("user_id")
	// Error handling jika userID tidak ditemukan atau tidak valid bisa ditambahkan di sini
	bag, err := ctl.service.CreateStickerBag(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, bag, "RackStagingSticker created", nil, http.StatusOK)
}

// List all rackStagingSticker
func (ctl *RackStagingStickerController) List(c *gin.Context) {
	bags, err := ctl.service.ListBags()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, bags, "List rackStagingSticker", nil, http.StatusOK)
}

// Get detail rackStagingSticker
func (ctl *RackStagingStickerController) GetDetail(c *gin.Context) {
	id := c.Param("id")
	bag, err := ctl.service.GetBagDetail(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, err.Error())
		return
	}
	utils.SendSuccess(c, bag, "Detail rackStagingSticker", nil, http.StatusOK)
}

// Scan product to rackStagingSticker (Bag) - hanya location staging_sticker
func (ctl *RackStagingStickerController) ScanBarcodeWarehouse(c *gin.Context) {
	bagID := c.Param("id")
	var req struct {
		BarcodeWarehouse string `json:"barcode_warehouse" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if bagID == "" {
		utils.SendError(c, http.StatusBadRequest, "Kolom bagID kosong")
		return
	}
	// Ambil product master by barcode_warehouse
	master, err := ctl.service.GetProductByBarcodeWarehouse(req.BarcodeWarehouse)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	if master.Location != "staging_sticker" {
		utils.SendError(c, http.StatusBadRequest, "Hanya produk dengan lokasi staging_sticker yang dapat di-scan ke rackStagingSticker")
		return
	}
	if master.BagID != nil && *master.BagID != "" {
		utils.SendError(c, http.StatusBadRequest, "Produk sudah dimasukkan ke bag, tidak dapat di scan ulang")
		return
	}
	// Assign product ke bag
	err = ctl.service.SetBag(master.ID.String(), bagID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Refresh master
	master, _ = ctl.service.GetProductByBarcodeWarehouse(req.BarcodeWarehouse)
	utils.SendSuccess(c, master, "Product assigned to rackStagingSticker", nil, http.StatusOK)
}

// List produk dalam bag (rackStagingSticker)
func (ctl *RackStagingStickerController) ListByBagID(c *gin.Context) {
	bagID := c.Param("id")
	if bagID == "" {
		utils.SendError(c, http.StatusBadRequest, "bagID is required")
		return
	}
	products, err := ctl.service.ListProductsByBagID(bagID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, products, "List produk dalam bag", nil, http.StatusOK)
}
