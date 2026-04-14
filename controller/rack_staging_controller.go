package controller

import (
	"net/http"

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
	utils.SendSuccess(ctx, result, "Rack staging created successfully", nil, http.StatusOK)
}

// Get detail of a rack staging
func (c *RackStagingController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.Service.GetRackStagingDetail(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.SendSuccess(ctx, result, "Rack staging detail retrieved successfully", nil, http.StatusOK)
}
