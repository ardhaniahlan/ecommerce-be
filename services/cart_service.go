package services

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/repositories"
)

type CartService interface {
	AddToCart(userID string, input dto.AddToCartInput) error
	GetCartItems(userID string) ([]dto.CartItemResponse, error)
	RemoveItem(cartID int, userID string) error
	UpdateItemQuantity(cartID int, userID string, quantity int) error
}

type cartService struct {
	repo repositories.CartRepository
}

func NewCartService(repo repositories.CartRepository) CartService {
	return &cartService{repo}
}

func (s *cartService) AddToCart(userID string, input dto.AddToCartInput) error {
	return s.repo.AddToCart(userID, input.ProductID, input.Quantity)
}

func (s *cartService) GetCartItems(userID string) ([]dto.CartItemResponse, error) {
	return s.repo.GetCartByUserID(userID)
}

func (s *cartService) RemoveItem(cartID int, userID string) error {
	return s.repo.RemoveFromCart(cartID, userID)
}

func (s *cartService) UpdateItemQuantity(cartID int, userID string, quantity int) error {
	return s.repo.UpdateCartItemQuantity(cartID, userID, quantity)
}