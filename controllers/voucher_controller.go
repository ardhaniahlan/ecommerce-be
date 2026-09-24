package controllers

import (
	"ecommerce-backend/models"
	"ecommerce-backend/services"
	"ecommerce-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoucherController struct {
	service services.VoucherService
}

func NewVoucherController(service services.VoucherService) *VoucherController {
	return &VoucherController{service}
}

func (c *VoucherController) CreateVoucher(ctx *gin.Context) {
	var input models.Voucher
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	if err := c.service.CreateVoucher(input); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menambahkan voucher", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Voucher berhasil ditambahkan", nil)
}

func (c *VoucherController) GetAllVoucher(ctx *gin.Context) {
	vouchers, err := c.service.GetAllVouchers()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data voucher", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Daftar voucher berhasil diambil", vouchers)
}

func (c *VoucherController) ApplyVoucher(ctx *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Kode promo tidak boleh kosong", err.Error())
		return
	}

	voucher, err := c.service.ApplyVoucher(req.Code)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Gagal menerapkan kode promo", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Kode promo berhasil diterapkan", gin.H{
		"code":                voucher.Code,
		"discount_percentage": voucher.DiscountPercentage,
		"max_discount_amount": voucher.MaxDiscountAmount,
	})
}
