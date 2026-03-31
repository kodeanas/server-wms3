package controller

import (
	"net/http"
	"strconv"

	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

// NewUserController creates a new user controller
func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} utils.Response
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.CreateUser(ctx.Request.Context(), &user); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SendSuccess(ctx, user, "User created successfully", http.StatusCreated)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.service.GetUser(ctx.Request.Context(), id)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccess(ctx, user, "User retrieved successfully")
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieve paginated list of users
// @Param limit query int false "Limit results" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.Response
// @Router /users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	limit := 10
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	users, total, err := c.service.ListUsers(ctx.Request.Context(), limit, offset)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"users": users,
			"total": total,
		},
		"message": "Users retrieved successfully",
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user
// @Param id path string true "User ID"
// @Param user body models.User true "User data"
// @Success 200 {object} utils.Response
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.UpdateUser(ctx.Request.Context(), id, &user); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.SendSuccess(ctx, user, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user from the system
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteUser(ctx.Request.Context(), id); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.SendSuccess(ctx, nil, "User deleted successfully")
}
