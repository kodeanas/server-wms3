package controller

import (
	"net/http"

	"wms/models"
	"wms/services"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

type RackDisplayController struct {
	Service *services.RackDisplayService
}

func NewRackDisplayController(service *services.RackDisplayService) *RackDisplayController {
	return &RackDisplayController{Service: service}
}

func (ctl *RackDisplayController) Create(c *gin.Context) {
	var rack models.RackDisplay
	if err := c.ShouldBindJSON(&rack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctl.Service.Create(&rack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rack)
}

func (ctl *RackDisplayController) GetAll(c *gin.Context) {
	racks, err := ctl.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, racks)
}

func (ctl *RackDisplayController) GetByID(c *gin.Context) {
	id := c.Param("id")
	rack, err := ctl.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rack)
}

func (ctl *RackDisplayController) Update(c *gin.Context) {
	id := c.Param("id")
	var rack models.RackDisplay
	if err := c.ShouldBindJSON(&rack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rack.ID = uuid.MustParse(id)
	if err := ctl.Service.Update(&rack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rack)
}

func (ctl *RackDisplayController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctl.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
