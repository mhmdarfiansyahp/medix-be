package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"medix-be/internal/transaction/model"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
	"time"
)

type TransactionService interface {
	CreateTransaction(userID uint, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetAllTransactions() ([]dto.TransactionResponse, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
	CancelTransaction(id uint, userID uint) error
	GetTodayTransactions(userID uint) (*dto.TodayTransactionResponse, error)
	GetReceipt(id uint) (*dto.ReceiptResponse, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(userID uint, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	var transaksi model.Transaksi
	var detailsEntities []model.DetailPembelian
	var totalHarga float64

	seenObat := make(map[uint]bool)

	err := s.repo.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Details {

			if seenObat[item.IDObat] {
			return errors.New(
				"the same medicine cannot be added more than once",
			)
			}

			seenObat[item.IDObat] = true

			// Fetch medicine from database
			obat, err := s.repo.FindObatByID(tx, item.IDObat)
			if err != nil {
				return errors.New("medicine not found")
			}

			// Ensure medicine is active
			if obat.Status != 1 {
			return fmt.Errorf(
				"medicine %s is currently inactive",
				obat.NamaObat,
			)
			}

			// Validate stock
			if obat.Stok < item.Jumlah {
			return fmt.Errorf(
				"insufficient stock for medicine %s",
				obat.NamaObat,
			)
			}

			// Price snapshot from database
			hargaSatuan := obat.Harga
			subtotal := float64(item.Jumlah) * hargaSatuan
			totalHarga += subtotal

			detailsEntities = append(
				detailsEntities,
				model.DetailPembelian{
					IDObat:      item.IDObat,
					Jumlah:      item.Jumlah,
					HargaSatuan: hargaSatuan,
					Subtotal:    subtotal,
				},
			)
		}

		// Buat transaksi utama
		transaksi = model.Transaksi{
			IDUser:     userID,
			TotalHarga: totalHarga,
			Status:     model.StatusTransaksiSelesai,
		}

		if err := s.repo.Create(tx, &transaksi); err != nil {
			return err
		}

		// Simpan detail dan kurangi stok
		for i := range detailsEntities {

			detailsEntities[i].IDTransaksi = transaksi.IDTransaksi

			// Simpan detail
			if err := s.repo.CreateDetail(tx, &detailsEntities[i]); err != nil {
				return err
			}

			if err := s.repo.UpdateStokObat(
				tx,
				detailsEntities[i].IDObat,
				detailsEntities[i].Jumlah,
			); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.New(
			"failed to process transaction: " + err.Error(),
		)
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
		return nil, errors.New("transaction not found")
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

func (s *transactionService) CancelTransaction(id uint, userID uint) error {
	tx := s.repo.GetDB().Begin()

	if tx.Error != nil {
		return tx.Error
	}

	var transaksi model.Transaksi

	err := tx.
		Preload("Details").
		Where("id_transaksi = ?", id).
		First(&transaksi).
		Error

	if err != nil {
		tx.Rollback()
		return errors.New("transaksi tidak ditemukan")
	}

	// Ensure transaction belongs to the logged-in cashier
	if transaksi.IDUser != userID {
		tx.Rollback()
		return errors.New("you do not have access to cancel this transaction")
	}

	// Already cancelled
	if transaksi.Status == model.StatusTransaksiDibatalkan {
		tx.Rollback()
		return errors.New("transaction already cancelled")
	}

	// Only today's transactions
	now := time.Now()

	if transaksi.TglTransaksi.Year() != now.Year() ||
		transaksi.TglTransaksi.YearDay() != now.YearDay() {
		tx.Rollback()
		return errors.New(
			"transaction can only be cancelled on the same day",
		)
	}

	// Restore stock
	for _, detail := range transaksi.Details {
		if err := s.repo.RestoreStokObat(
			tx,
			detail.IDObat,
			detail.Jumlah,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Change status to cancelled
	if err := s.repo.UpdateStatus(tx, id, model.StatusTransaksiDibatalkan); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *transactionService) GetTodayTransactions(userID uint) (*dto.TodayTransactionResponse, error) {

	transactions, err := s.repo.FindTodayByUser(userID)
	if err != nil {
		return nil, err
	}

	totalPenjualan, totalTransaksi, err :=
		s.repo.GetTodaySummary(userID)

	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		result = append(result, toTransactionResponse(*transaction))
	}

	return &dto.TodayTransactionResponse{
		Transactions: result,
		Summary: dto.TransactionSummaryResponse{
			TotalTransaksi: totalTransaksi,
			TotalPenjualan: totalPenjualan,
		},
	}, nil
}

func (s *transactionService) GetReceipt(id uint) (*dto.ReceiptResponse, error) {

	transaction, err := s.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if transaction.Status == model.StatusTransaksiDibatalkan {
		return nil, errors.New("transaksi sudah dibatalkan")
	}

	details := make([]dto.DetailItemResponse, 0)

	for _, detail := range transaction.Details {
		details = append(details, dto.DetailItemResponse{
			IDDetail:    detail.IDDetail,
			IDObat:      detail.IDObat,
			Jumlah:      detail.Jumlah,
			HargaSatuan: detail.HargaSatuan,
			Subtotal:    detail.Subtotal,
		})
	}

	return &dto.ReceiptResponse{
		IDTransaksi:  transaction.IDTransaksi,
		TglTransaksi: transaction.TglTransaksi,
		IDUser:       transaction.IDUser,
		Details:      details,
		TotalHarga:   transaction.TotalHarga,
	}, nil
}
