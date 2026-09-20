package repositories

import (
	"ecommerce-backend/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VoucherRepository interface {
	CreateVoucher(v models.Voucher) error
	GetAllVouchers() ([]models.Voucher, error)
}

type voucherRepository struct {
	db *sqlx.DB
}

func NewVoucherRepository(db *sqlx.DB) VoucherRepository {
	return &voucherRepository{db}
}


func (r *voucherRepository) CreateVoucher(v models.Voucher) error {
	query := `
		INSERT INTO vouchers (id, code, discount_percentage, max_discount_amount, quota, valid_from, valid_until)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	newID := "vch-" + uuid.New().String()
	
	_, err := r.db.Exec(query, newID, v.Code, v.DiscountPercentage, v.MaxDiscountAmount, v.Quota, v.ValidFrom, v.ValidUntil)
	return err
}

func (r *voucherRepository) GetAllVouchers() ([]models.Voucher, error) {
	var vouchers []models.Voucher
	err := r.db.Select(&vouchers, "SELECT id, code, discount_percentage, max_discount_amount, quota, valid_from, valid_until FROM vouchers ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	if vouchers == nil {
		vouchers = []models.Voucher{}
	}
	return vouchers, nil
}
