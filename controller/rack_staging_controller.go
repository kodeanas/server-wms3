package controller

import (
	"net/http"

	"wms/config"
	"wms/repositories"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type RackStagingController struct {
	Service *services.RackStagingService
}

func NewRackStagingController(service *services.RackStagingService) *RackStagingController {
	return &RackStagingController{Service: service}
}

func (c *RackStagingController) Create(ctx *gin.Context) {
	var req struct {
		RackDisplayID string `json:"rack_display_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.Service.CreateRackStaging(req.RackDisplayID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.SendSuccess(ctx, result, "Rack staging dibuat berhasil", nil, http.StatusOK)
}

// Get detail of a rack staging
func (c *RackStagingController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("rackStagingID")
	result, err := c.Service.GetRackStagingDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.SendSuccess(ctx, result, "Rack staging detail didapatkan", nil, http.StatusOK)
}

// List all rack stagings
func (c *RackStagingController) ListAll(ctx *gin.Context) {
	racks, err := c.Service.ListAllRackStaging()
	if err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccess(ctx, racks, "List semua rack stagings", nil, http.StatusOK)
}

// Finish rack staging: set is_moved dan update semua product master ke display
func (c *RackStagingController) Finish(ctx *gin.Context) {
	rackStagingID := ctx.Param("rackStagingID")
	if rackStagingID == "" {
		utils.SendError(ctx, 400, "Kolom rackStagingID kosong")
		return
	}
	// Langsung pakai repository, tidak perlu type assertion interface
	productMasterRepo := repositories.NewProductMasterRepository(config.DB)
	err := c.Service.FinishRackStaging(rackStagingID, productMasterRepo)
	if err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccess(ctx, nil, "Rack staging selesai dan produk pindah ke display", nil, http.StatusOK)
}
