package services

import "ecommerce-backend/repositories"

type OrderService interface {
	CheckoutCart(userID string) (string, error)
}

type orderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{repo}
}

func (s *orderService) CheckoutCart(userID string) (string, error) {
	return s.repo.Checkout(userID)
}