package controller

import (
	"net/http"
	dto "wms/dto/response"
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

	utils.SendItemSuccess(c, sticker, "", http.StatusOK)
}

// ListStickers endpoint.
func (ctrl *StickerController) ListStickers(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	stickers, total, err := ctrl.service.ListStickersPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}

	resp := make([]dto.StickerResponse, 0, len(stickers))

	for _, s := range stickers {
		resp = append(resp, dto.StickerResponse{
			ID:         s.ID.String(),
			CodeHex:    s.CodeHex,
			Name:       s.Name,
			Slug:       s.Slug,
			Type:       s.Type,
			FixedPrice: s.FixedPrice,
			MinPrice:   (*float64)(s.MinPrice),
			MaxPrice:   (*float64)(s.MaxPrice),
			Status:     s.Status,
			CreatedAt:  s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			DeletedAt:  s.DeletedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	utils.SendListSuccess(c, stickers, pg.Page, pg.Limit, total, "", http.StatusOK)
}

func (ctrl *StickerController) ListStickerSelect(c *gin.Context) {
	search := c.Query("search")

	stickers, err := ctrl.service.ListStickers(search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}

	resp := make([]dto.ListStickerSelect, 0, len(stickers))
	for _, s := range stickers {
		resp = append(resp, dto.ListStickerSelect{
			ID:         s.ID.String(),
			CodeHexx:   s.CodeHex,
			Name:       s.Name,
			Slug:       s.Slug,
			Type:       s.Type,
			FixedPrice: s.FixedPrice,
			MinPrice:   (*float64)(s.MaxPrice),
			MaxPrice:   (*float64)(s.MaxPrice),
			Status:     s.Status,
		})
	}

	utils.SendItemSuccess(c, resp, "List sticker berhasil di ambil", http.StatusOK)
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
