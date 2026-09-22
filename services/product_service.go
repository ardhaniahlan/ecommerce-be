package services

import (
	"context"
	"ecommerce-backend/dto"
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
	"ecommerce-backend/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductService interface {
	CreateProduct(input dto.CreateProductInput) (models.Product, error)
	GetAllProducts(page int, limit int, search string) ([]models.Product, utils.Pagination, error)
	GetProductByID(id int) (models.Product, error)
	UpdateProduct(id int, input dto.CreateProductInput) (models.Product, error)
	DeleteProduct(id int) error
}

type productService struct {
	repo  repositories.ProductRepository
	redis *redis.Client
}

func NewProductService(repo repositories.ProductRepository, redisClient *redis.Client) ProductService {
	return &productService{repo, redisClient}
}

type PaginatedProducts struct {
	Products []models.Product `json:"products"`
	Meta     utils.Pagination `json:"meta"`
}

func (s *productService) CreateProduct(input dto.CreateProductInput) (models.Product, error) {
	product := models.Product{
		Name:               input.Name,
		Description:        input.Description,
		Price:              input.Price,
		Stock:              input.Stock,
		DiscountPercentage: input.DiscountPercentage,
		DiscountStart:      input.DiscountStart,
		DiscountEnd:        input.DiscountEnd,
	}

	err := s.repo.Create(&product, input.ImageURLs)

	if err == nil {
		cacheKey := fmt.Sprintf("products:page:%d:limit:%d:search:%s", 1, 10, "")
		s.redis.Del(context.Background(), cacheKey)
	}
	return product, err
}

func (s *productService) GetAllProducts(page, limit int, search string) ([]models.Product, utils.Pagination, error) {
	cacheKey := fmt.Sprintf("products:page:%d:limit:%d:search:%s", page, limit, search)

	cachedData, err := utils.GetOrSetCache(context.Background(), s.redis, cacheKey, 2*time.Minute, func() (PaginatedProducts, error) {

		products, totalItems, errRepo := s.repo.GetAll(page, limit, search)
		if errRepo != nil {
			return PaginatedProducts{}, errRepo
		}

		meta := utils.GeneratePagination(page, limit, totalItems)
		return PaginatedProducts{Products: products, Meta: meta}, nil
	})

	if err != nil {
		return nil, utils.Pagination{}, err
	}

	return cachedData.Products, cachedData.Meta, nil
}

func (s *productService) GetProductByID(id int) (models.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	return utils.GetOrSetCache(context.Background(), s.redis, cacheKey, 5*time.Minute, func() (models.Product, error) {
		return s.repo.GetByID(id)
	})
}

func (s *productService) UpdateProduct(id int, input dto.CreateProductInput) (models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return product, err
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.DiscountPercentage = input.DiscountPercentage
	product.DiscountStart = input.DiscountStart
	product.DiscountEnd = input.DiscountEnd

	err = s.repo.Update(&product, input.ImageURLs)

	if err == nil {
		cacheKey := fmt.Sprintf("product:%d", id)
		s.redis.Del(context.Background(), cacheKey)
	}

	return product, err
}

func (s *productService) DeleteProduct(id int) error {
	err := s.repo.Delete(id)

	if err == nil {
		cacheKey := fmt.Sprintf("product:%d", id)
		s.redis.Del(context.Background(), cacheKey)
	}

	return err
}
