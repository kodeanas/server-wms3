package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type WholesaleBagController struct {
	service services.WholesaleBagService
}

func NewWholesaleBagController(service services.WholesaleBagService) *WholesaleBagController {
	return &WholesaleBagController{service: service}
}

func (ctl *WholesaleBagController) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	bag, err := ctl.service.CreateWholesaleBag(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(c, bag, "WholesaleBag created", nil, http.StatusOK)
}

func (ctl *WholesaleBagController) List(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	bags, total, err := ctl.service.ListWholesaleBagsPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendListSuccess(c, bags, pg.Page, pg.Limit, total, "", http.StatusOK)
}

func (ctl *WholesaleBagController) GetDetail(c *gin.Context) {
	bagID := c.Param("bagID")
	detail, err := ctl.service.GetWholesaleBagDetail(bagID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendItemSuccess(c, detail, "", http.StatusOK)
}

func (ctl *WholesaleBagController) ListByBagID(c *gin.Context) {
	bagID := c.Param("bagID")
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	products, total, err := ctl.service.ListProductsByWholesaleBagIDPaginated(bagID, pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendListSuccess(c, products, pg.Page, pg.Limit, total, "", http.StatusOK)
}

func (ctl *WholesaleBagController) ScanBarcodeWarehouse(c *gin.Context) {
	var req struct {
		BarcodeWarehouse string `json:"barcode_warehouse" binding:"required"`
	}

	bagID := c.Param("bagID")
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if bagID == "" {
		utils.SendError(c, http.StatusBadRequest, "Kolom bagID kosong")
		return
	}
	master, err := ctl.service.GetProductByBarcodeWarehouse(req.BarcodeWarehouse)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}
	if master.Location != "staging_reguler" && master.Location != "display" {
		utils.SendError(c, http.StatusBadRequest, "Hanya produk dengan lokasi staging_reguler atau display yang dapat di-scan ke wholesale bag")
		return
	}
	if master.BagID != nil && *master.BagID != "" {
		utils.SendError(c, http.StatusBadRequest, "Produk sudah dimasukkan ke bag, tidak dapat di scan ulang")
		return
	}
	err = ctl.service.SetBag(master.ID.String(), bagID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	master, _ = ctl.service.GetProductByBarcodeWarehouse(req.BarcodeWarehouse)
	utils.SendSuccess(c, master, "Product assigned to wholesale bag", nil, http.StatusOK)
}
