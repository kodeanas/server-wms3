package controller

import (
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaxController struct {
	service *services.TaxService
}

func NewTaxController(service *services.TaxService) *TaxController {
	return &TaxController{service: service}
}

func (c *TaxController) Create(ctx *gin.Context) {
	var tax models.Tax
	if err := ctx.ShouldBindJSON(&tax); err != nil {
		errors := []utils.ErrorItem{{Field: "body", Message: err.Error()}}
		utils.SendValidationError(ctx, errors)
		return
	}
	if err := c.service.Create(&tax); err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, tax, "Berhasil membuat tax", 201)
}

func (c *TaxController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var tax models.Tax
	if err := ctx.ShouldBindJSON(&tax); err != nil {
		errors := []utils.ErrorItem{{Field: "body", Message: err.Error()}}
		utils.SendValidationError(ctx, errors)
		return
	}
	tax.ID = uuid.MustParse(id)
	if err := c.service.Update(&tax); err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, tax, "Berhasil update tax")
}

func (c *TaxController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, nil, "Berhasil hapus tax")
}

func (c *TaxController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	tax, err := c.service.FindByID(id)
	if err != nil {
		utils.SendError(ctx, 404, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, tax, "Berhasil ambil detail tax")
}

func (c *TaxController) List(ctx *gin.Context) {
	taxes, err := c.service.FindAll()
	if err != nil {
		utils.SendError(ctx, 500, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, taxes, "Berhasil ambil list tax")
}

func (c *TaxController) GetActive(ctx *gin.Context) {
	tax, err := c.service.FindActive()
	if err != nil {
		utils.SendError(ctx, 404, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(ctx, tax, "Berhasil ambil tax aktif")
}
