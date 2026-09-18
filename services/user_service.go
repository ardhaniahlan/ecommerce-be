package services

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
)

type UserService interface {
	UpdateUserProfile(userID string, req dto.UpdateProfileRequest) error
	GetUserProfile(userID string) (models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) UpdateUserProfile(userID string, req dto.UpdateProfileRequest) error {
	return s.userRepo.UpdateProfile(userID, req.FullName, req.Phone, req.Province, req.City, req.District, req.PostalCode, req.StreetAddress)
}

func (s *userService) GetUserProfile(userID string) (models.User, error) {
	return s.userRepo.GetProfile(userID)
}
