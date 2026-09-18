package repositories

import (
	"ecommerce-backend/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	UpdateProfile(userID, name, phone, province, city, district, postalCode, streetAddress string) error
	GetProfile(userID string) (models.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) UpdateProfile(userID, fullName, phone, province, city, district, postalCode, streetAddress string) error {
	query := `UPDATE users 
	          SET full_name = $1, phone = $2, province = $3, city = $4, district = $5, postal_code = $6, street_address = $7 
	          WHERE id = $8`
	_, err := r.db.Exec(query, fullName, phone, province, city, district, postalCode, streetAddress, userID)
	return err
}

func (r *userRepository) GetProfile(userID string) (models.User, error) {
	var profile models.User
	query := `
		SELECT id, full_name, email, phone, role, province, city, district, postal_code, street_address 
		FROM users WHERE id = $1
	`
	err := r.db.Get(&profile, query, userID)
	return profile, err
}