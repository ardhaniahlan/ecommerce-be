package controllers

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/services"
	"ecommerce-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data profil tidak lengkap", err.Error())
		return
	}

	if err := c.service.UpdateUserProfile(userID, req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal memperbarui profil", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profil dan alamat berhasil diperbarui", nil)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	profile, err := c.service.GetUserProfile(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "User tidak ditemukan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Berhasil mengambil data profil", profile)
}