package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var payload services.CreateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	user, err := ctrl.service.CreateUser(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, user, "User berhasil ditambahkan", nil, http.StatusCreated)
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := ctrl.service.GetUserByID(id)
	if err != nil {
		utils.SendError(c, 404, err.Error())
		return
	}
	utils.SendSuccess(c, user, "OK", nil, http.StatusOK)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	users, err := ctrl.service.ListUsers()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, users, "OK", nil, http.StatusOK)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var payload services.UpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	user, err := ctrl.service.UpdateUser(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, user, "User berhasil diupdate", nil, http.StatusOK)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteUser(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "User berhasil dihapus", nil, http.StatusOK)
}

func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	type PasswordPayload struct {
		Password string `json:"password" binding:"required"`
	}
	var payload PasswordPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	if err := ctrl.service.UpdatePassword(id, payload.Password); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Password berhasil diupdate", nil, http.StatusOK)
}
