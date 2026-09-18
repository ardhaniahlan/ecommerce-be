package controllers

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/services"
	"ecommerce-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var input dto.CreateProductInput
	
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	product, err := c.service.CreateProduct(input)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menambahkan produk", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Produk berhasil ditambahkan", product)
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	search := ctx.Query("search")


	products, meta, err := c.service.GetAllProducts(page, limit, search)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengambil data produk", err.Error())
		return
	}

	utils.SuccessResponseWithMeta(ctx, http.StatusOK, "Daftar produk berhasil diambil", products, meta)
}

func (c *ProductController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	product, err := c.service.GetProductByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Produk tidak ditemukan", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Detail produk berhasil diambil", product)
}

func (c *ProductController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var input dto.CreateProductInput
	
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	product, err := c.service.UpdateProduct(id, input)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal mengupdate produk", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Produk berhasil diupdate", product)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.service.DeleteProduct(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Gagal menghapus produk", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Produk berhasil dihapus (Soft Delete)", nil)
}