package controllers

import (
	"ecommerce-backend/services"
	"ecommerce-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service services.OrderService
}

func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{service}
}

func (c *OrderController) Checkout(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	orderID, err := c.service.CheckoutCart(userID)
	if err != nil {
		if err.Error() == "keranjang belanja kosong" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Checkout gagal", err.Error())
			return
		}
		
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses pesanan", err.Error())
		return
	}
	responseData := map[string]string{
		"orderId": orderID,
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Pesanan berhasil dibuat", responseData)
}