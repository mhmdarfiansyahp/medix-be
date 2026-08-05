package service

import (
	"errors"
	"gorm.io/gorm"
	"medix-be/internal/transaction/model"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
)

type TransactionService interface {
	CreateTransaction(req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetAllTransactions() ([]dto.TransactionResponse, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	var totalHarga float64
	var detailsEntities []model.DetailPembelian

	// Hitung total harga & siapkan data detail
	for _, item := range req.Details {
		subtotal := float64(item.Jumlah) * item.HargaSatuan
		totalHarga += subtotal

		detailsEntities = append(detailsEntities, model.DetailPembelian{
			IDObat:      item.IDObat,
			Jumlah:      item.Jumlah,
			HargaSatuan: item.HargaSatuan,
			Subtotal:    subtotal,
		})
	}

	transaksi := model.Transaksi{
		IDUser:     req.IDUser,
		TotalHarga: totalHarga,
		Status:     "selesai",
	}

	// Jalankan transaksi DB (ACID safe)
	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Header Transaksi
		if err := s.repo.Create(tx, &transaksi); err != nil {
			return err
		}

		// 2. Simpan Detail & Potong Stok Obat
		for i := range detailsEntities {
			detailsEntities[i].IDTransaksi = transaksi.IDTransaksi
			if err := s.repo.CreateDetail(tx, &detailsEntities[i]); err != nil {
				return err
			}

			// Potong stok obat
			err := s.repo.UpdateStokObat(tx, detailsEntities[i].IDObat, detailsEntities[i].Jumlah)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.New("gagal memproses transaksi: " + err.Error())
	}

	transaksi.Details = detailsEntities
	res := toTransactionResponse(transaksi)
	return &res, nil
}

func (s *transactionService) GetAllTransactions() ([]dto.TransactionResponse, error) {
	transaksis, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.TransactionResponse
	for _, t := range transaksis {
		responses = append(responses, toTransactionResponse(t))
	}
	return responses, nil
}

func (s *transactionService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	transaksi, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	res := toTransactionResponse(*transaksi)
	return &res, nil
}

func toTransactionResponse(t model.Transaksi) dto.TransactionResponse {
	var detailsRes []dto.DetailItemResponse
	for _, d := range t.Details {
		detailsRes = append(detailsRes, dto.DetailItemResponse{
			IDDetail:    d.IDDetail,
			IDObat:      d.IDObat,
			Jumlah:      d.Jumlah,
			HargaSatuan: d.HargaSatuan,
			Subtotal:    d.Subtotal,
		})
	}

	return dto.TransactionResponse{
		IDTransaksi:  t.IDTransaksi,
		IDUser:       t.IDUser,
		TglTransaksi: t.TglTransaksi,
		TotalHarga:   t.TotalHarga,
		Status:       t.Status,
		Details:      detailsRes,
	}
}
