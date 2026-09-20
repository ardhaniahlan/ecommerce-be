package controllers

import (
	"ecommerce-backend/dto"
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

	var req dto.CheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.VoucherCode = ""
	}

	result, err := c.service.CheckoutCart(userID, req.VoucherCode)
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

func (c *OrderController) GetAllOrdersAdmin(ctx *gin.Context) {
	orders, err := c.service.GetAllOrdersAdmin()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil daftar pesanan", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Berhasil mengambil semua pesanan", orders)
}

func (c *OrderController) InputTracking(ctx *gin.Context) {
	orderID := ctx.Param("id")

	var req dto.TrackingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Nomor resi tidak valid", err.Error())
		return
	}

	if err := c.service.InputTrackingNumber(orderID, req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal mengupdate resi", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Nomor resi berhasil disimpan dan status menjadi Shipped", nil)
}

func (c *OrderController) GetUserHistory(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)

	orders, err := c.service.GetUserOrderHistory(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil riwayat pesanan", err.Error())
		return
	}

	if orders == nil {
		orders = []dto.AdminOrderResponse{}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Berhasil mengambil riwayat pesanan", orders)
}

func (c *OrderController) CompleteOrder(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(string)
	orderID := ctx.Param("id")

	if err := c.service.CompleteOrder(orderID, userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal menyelesaikan pesanan", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Pesanan berhasil diselesaikan", nil)
}

func (c *OrderController) GetStats(ctx *gin.Context) {
	stats, err := c.service.GetAdminStats()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil statistik dashboard", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Berhasil mengambil statistik", stats)
}