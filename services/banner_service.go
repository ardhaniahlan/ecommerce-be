package services

import (
	"ecommerce-backend/dto"
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
)

type BannerService interface {
	CreateBanner(input dto.BannerInput) (error)
	GetActiveBanners() ([]models.Banner, error)
	DeleteBanner(id int) error
}

type bannerService struct {
	repo repositories.BannerRepository
}

func NewBannerService(repo repositories.BannerRepository) BannerService {
	return &bannerService{repo}
}

func (s *bannerService) CreateBanner(input dto.BannerInput) (error) {
	banner := models.Banner{
		Title: input.Title,
		ImageURL: input.ImageURL,
	}
	return s.repo.Create(&banner)
}

func (s *bannerService) GetActiveBanners() ([]models.Banner, error) {
	return s.repo.GetActiveBanners()
}

func (s *bannerService) DeleteBanner(id int) error {
	return s.repo.Delete(id)
}
