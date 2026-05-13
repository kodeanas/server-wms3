package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type CargoController struct {
	service services.CargoService
}

func NewCargoController(service services.CargoService) *CargoController {
	return &CargoController{service: service}
}

func (ctl *CargoController) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	cargo, err := ctl.service.CreateCargo(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, cargo, "Cargo created", nil, http.StatusOK)
}

func (ctl *CargoController) List(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	cargos, total, err := ctl.service.ListCargosPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendListSuccess(c, cargos, pg.Page, pg.Limit, total, "", http.StatusOK)
}

func (ctl *CargoController) GetDetail(c *gin.Context) {
	cargoID := c.Param("id")
	detail, err := ctl.service.GetCargoDetail(cargoID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendItemSuccess(c, detail, "", http.StatusOK)
}

func (ctl *CargoController) ScanBag(c *gin.Context) {
	var req struct {
		BagCode string `json:"bag_code" binding:"required"`
	}

	cargoID := c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctl.service.ScanBag(cargoID, req.BagCode); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, nil, "Bag scanned successfully", nil, http.StatusOK)
}

func (ctl *CargoController) Finish(c *gin.Context) {
	cargoID := c.Param("id")

	if err := ctl.service.FinishCargo(cargoID); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccess(c, nil, "Cargo finished", nil, http.StatusOK)
}
