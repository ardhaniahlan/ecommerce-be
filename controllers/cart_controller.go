package controllers

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/repositories"
	"ecommerce-backend/services"
	"ecommerce-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	service services.CartService
}

func NewCartController(service services.CartService) *CartController {
	return &CartController{service}
}

func (c *CartController) AddItem(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	var input dto.AddToCartInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	if err := c.service.AddToCart(userID, input); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menambahkan ke keranjang", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Barang berhasil ditambahkan ke keranjang", nil)
}

func (c *CartController) GetCart(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	items, err := c.service.GetCartItems(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data keranjang", err.Error())
		return
	}

	if items == nil {
		items = []repositories.CartItemResponse{}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Data keranjang berhasil diambil", items)
}

func (c *CartController) RemoveItem(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	cartID, _ := strconv.Atoi(ctx.Param("id"))

	if err := c.service.RemoveItem(cartID, userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus barang dari keranjang", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Barang dihapus dari keranjang", nil)
}

func (c *CartController) UpdateQuantity(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	
	cartID, _ := strconv.Atoi(ctx.Param("id"))

	var input dto.UpdateCartInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	if err := c.service.UpdateItemQuantity(cartID, userID, input.Quantity); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "barang tidak ditemukan di keranjang Anda" {
			statusCode = http.StatusNotFound
		}
		
		utils.ErrorResponse(ctx, statusCode, "Gagal mengupdate jumlah barang", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Jumlah barang berhasil diupdate", nil)
}