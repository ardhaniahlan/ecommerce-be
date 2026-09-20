package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
)

type VoucherService interface {
	CreateVoucher(voucher models.Voucher) error
	GetAllVouchers() ([]models.Voucher, error)
}

type voucherService struct {
	voucherRepo repositories.VoucherRepository
}

func NewVoucherService(voucherRepo repositories.VoucherRepository) VoucherService {
	return &voucherService{voucherRepo}
}


func (s *voucherService) CreateVoucher(voucher models.Voucher) error {
	return s.voucherRepo.CreateVoucher(voucher)
}


func (s *voucherService) GetAllVouchers() ([]models.Voucher, error) {
	return s.voucherRepo.GetAllVouchers()
}