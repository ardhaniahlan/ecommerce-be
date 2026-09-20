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