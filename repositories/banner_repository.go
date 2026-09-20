package repositories

import (
	"ecommerce-backend/models"
	"github.com/jmoiron/sqlx"
)

type BannerRepository interface {
	Create(banner *models.Banner) error
	GetActiveBanners() ([]models.Banner, error)
	Delete(id int) error
}

type bannerRepository struct {
	db *sqlx.DB
}

func NewBannerRepository(db *sqlx.DB) BannerRepository {
	return &bannerRepository{db}
}

func (r *bannerRepository) Create(banner *models.Banner) error {
	query := `
		INSERT INTO banners (title, image_url) 
		VALUES (:title, :image_url) 
		RETURNING id, is_active, created_at`
	
	rows, err := r.db.NamedQuery(query, banner)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(banner)
	}
	return err
}

func (r *bannerRepository) GetActiveBanners() ([]models.Banner, error) {
	var banners []models.Banner
	query := `SELECT id, title, image_url, is_active, created_at FROM banners WHERE is_active = true ORDER BY id DESC`
	err := r.db.Select(&banners, query)
	return banners, err
}

func (r *bannerRepository) Delete(id int) error {
	query := `UPDATE banners SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}