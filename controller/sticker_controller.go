package controller

import (
	"net/http"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

// StickerController defines handlers for sticker resources.
type StickerController struct {
	service services.StickerService
}

// NewStickerController constructor.
func NewStickerController(service services.StickerService) *StickerController {
	return &StickerController{service: service}
}

// CreateSticker endpoint.
func (ctrl *StickerController) CreateSticker(c *gin.Context) {
	var payload models.CreateStickerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	sticker, err := ctrl.service.CreateSticker(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, sticker, "Sticker berhasil ditambahkan", nil, http.StatusCreated)
}

// GetStickerByID endpoint.
func (ctrl *StickerController) GetStickerByID(c *gin.Context) {
	id := c.Param("id")

	sticker, err := ctrl.service.GetStickerByID(id)
	if err != nil {
		utils.SendError(c, 404, "Sticker tidak ditemukan")
		return
	}

	utils.SendSuccess(c, sticker, "Sticker ditemukan", nil, http.StatusOK)
}

// ListStickers endpoint.
func (ctrl *StickerController) ListStickers(c *gin.Context) {
	stickers, err := ctrl.service.ListStickers()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}

	utils.SendSuccess(c, stickers, "Daftar sticker", nil, http.StatusOK)
}

// UpdateSticker endpoint.
func (ctrl *StickerController) UpdateSticker(c *gin.Context) {
	id := c.Param("id")

	var payload models.UpdateStickerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	sticker, err := ctrl.service.UpdateSticker(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, sticker, "Sticker berhasil diperbarui", nil, http.StatusOK)
}

// DeleteSticker endpoint.
func (ctrl *StickerController) DeleteSticker(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.service.DeleteSticker(id)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, nil, "Sticker berhasil dihapus", nil, http.StatusOK)
}
