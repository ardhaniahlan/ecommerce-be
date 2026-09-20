package services

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
	"ecommerce-backend/utils"
)

type ProductService interface {
	CreateProduct(input dto.CreateProductInput) (models.Product, error)
	GetAllProducts(page int, limit int, search string) ([]models.Product, utils.Pagination, error)
	GetProductByID(id int) (models.Product, error)
	UpdateProduct(id int, input dto.CreateProductInput) (models.Product, error)
	DeleteProduct(id int) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) CreateProduct(input dto.CreateProductInput) (models.Product, error) {
	product := models.Product{
		Name:               input.Name,
		Description:        input.Description,
		Price:              input.Price,
		Stock:              input.Stock,
		ImageURL:           input.ImageURL,
		DiscountPercentage: input.DiscountPercentage,
		DiscountStart:      input.DiscountStart,
		DiscountEnd:        input.DiscountEnd,
	}

	err := s.repo.Create(&product)
	return product, err
}

func (s *productService) GetAllProducts(page, limit int, search string) ([]models.Product, utils.Pagination, error) {
	products, totalItems, err := s.repo.GetAll(page, limit, search)
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	meta := utils.GeneratePagination(page, limit, totalItems)
	return products, meta, nil
}

func (s *productService) GetProductByID(id int) (models.Product, error) {
	return s.repo.GetByID(id)
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
	product.ImageURL = input.ImageURL
	product.DiscountPercentage = input.DiscountPercentage
	product.DiscountStart = input.DiscountStart
	product.DiscountEnd = input.DiscountEnd

	err = s.repo.Update(&product)
	return product, err
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
