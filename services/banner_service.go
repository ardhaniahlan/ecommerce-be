package services

import (
	"context"
	"ecommerce-backend/dto"
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
	"ecommerce-backend/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type BannerService interface {
	CreateBanner(input dto.BannerInput) (error)
	GetActiveBanners() ([]models.Banner, error)
	DeleteBanner(id int) error
}

type bannerService struct {
	repo repositories.BannerRepository
	redis *redis.Client
}	

func NewBannerService(repo repositories.BannerRepository, redisClient *redis.Client) BannerService {
	return &bannerService{repo, redisClient}
}

func (s *bannerService) CreateBanner(input dto.BannerInput) (error) {
	banner := models.Banner{
		Title: input.Title,
		ImageURL: input.ImageURL,
	}

	err := s.repo.Create(&banner)
	if err == nil {
		s.redis.Del(context.Background(), "banner:active")
	}
	return err
}

func (s *bannerService) GetActiveBanners() ([]models.Banner, error) {
	cacheKey := "banner:active"
	return utils.GetOrSetCache(context.Background(), s.redis, cacheKey, 5*time.Minute, func() ([]models.Banner, error) {
		return s.repo.GetActiveBanners()
	})
}

func (s *bannerService) DeleteBanner(id int) error {
	err := s.repo.Delete(id)
	if err == nil {
		cacheKey := "banner:active"
		s.redis.Del(context.Background(), cacheKey)
	}
	return err
}
