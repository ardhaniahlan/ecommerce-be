package services

import (
	"ecommerce-backend/models"
	"ecommerce-backend/repositories"
	"errors"
	"time"
)

type VoucherService interface {
	CreateVoucher(voucher models.Voucher) error
	GetAllVouchers() ([]models.Voucher, error)
	ApplyVoucher(code string) (*models.Voucher, error)
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

func (s *voucherService) ApplyVoucher(code string) (*models.Voucher, error) {
	voucher, err := s.voucherRepo.GetVoucherByCode(code)
	if err != nil {
		return nil, errors.New("kode promo tidak ditemukan atau tidak valid")
	}

	if voucher.Quota <= 0 {
		return nil, errors.New("kuota kode promo ini sudah habis")
	}

	now := time.Now()
	if voucher.ValidFrom == nil || now.Before(*voucher.ValidFrom) {
		return nil, errors.New("kode promo belum aktif")
	}
	if voucher.ValidUntil == nil || now.After(*voucher.ValidUntil) {
		return nil, errors.New("kode promo sudah kadaluarsa")
	}

	return &voucher, nil
}