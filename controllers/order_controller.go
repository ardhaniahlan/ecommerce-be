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

	result, err := c.service.CheckoutCart(userID)
	if err != nil {
		if err.Error() == "keranjang belanja kosong" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Checkout gagal", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal memproses pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Pesanan berhasil dibuat", result)
}

func (c *OrderController) Webhook(ctx *gin.Context) {
	var payload map[string]interface{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	err := c.service.ProcessMidtransWebhook(payload)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
