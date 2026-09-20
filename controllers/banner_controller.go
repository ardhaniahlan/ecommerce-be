package controllers

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BannerController struct {
	service services.BannerService
}

func NewBannerController(service services.BannerService) *BannerController {
	return &BannerController{service}
}

func (c *BannerController) CreateBanner(ctx *gin.Context) {
	var input dto.BannerInput
	
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid"})
		return
	}

	if err := c.service.CreateBanner(input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Banner berhasil dibuat"})
}

func (c *BannerController) GetActiveBanners(ctx *gin.Context) {
	banners, err := c.service.GetActiveBanners()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil banner",
		"data": banners,
	})
}

func (c *BannerController) DeleteBanner(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := c.service.DeleteBanner(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Banner berhasil dihapus"})
}