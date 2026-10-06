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
	AddToCart(userID uint, req dto.AddToCartRequest) (*dto.TransactionResponse, error)
	GetAllTransactions() ([]dto.TransactionResponse, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
CancelTransaction(id uint, userID uint, userRole string) error
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

func (s *transactionService) AddToCart(userID uint, req dto.AddToCartRequest) (*dto.TransactionResponse, error) {
	var transaksi model.Transaksi

	// Find today's transaction for this user (if any)
	var existingTransaksi model.Transaksi
	error := s.repo.GetDB().
		Preload("Details").
		Where("id_user = ? AND status = ?", userID, model.StatusTransaksiSelesai).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		First(&existingTransaksi).
		Error

	if error == nil && !errors.Is(error, gorm.ErrRecordNotFound) {
			// Found existing transaction
			transaksi = existingTransaksi
		} else {
			// Create new transaction
			transaksi = model.Transaksi{
				IDUser:  userID,
				Status:  model.StatusTransaksiSelesai,
			}
			if err := s.repo.Create(s.repo.GetDB(), &transaksi); err != nil {
				return nil, fmt.Errorf("failed to create transaction: %w", err)
			}
		}

	// Add medicine to transaction
	obat, err := s.repo.FindObatByID(s.repo.GetDB(), req.IDObat)
	if err != nil {
		return nil, fmt.Errorf("medicine not found: %w", err)
	}

	// Ensure medicine is active
	if obat.Status != 1 {
		return nil, fmt.Errorf("medicine %s is currently inactive", obat.NamaObat)
	}

	// Check if medicine already in this transaction
	for _, detail := range transaksi.Details {
		if detail.IDObat == req.IDObat {
			return nil, errors.New("the same medicine cannot be added more than once")
		}
	}

	// Check stock
	if obat.Stok < req.Jumlah {
		return nil, fmt.Errorf("insufficient stock for medicine %s", obat.NamaObat)
	}

	// Calculate price
	hargaSatuan := obat.Harga
	subtotal := float64(req.Jumlah) * hargaSatuan
	transaksi.TotalHarga += subtotal

	// Update transaction total
	if err := s.repo.UpdateStatus(s.repo.GetDB(), transaksi.IDTransaksi, model.StatusTransaksiSelesai); err != nil {
		return nil, fmt.Errorf("failed to update transaction total: %w", err)
	}

	// Create detail entity
	detailEntity := model.DetailPembelian{
		IDTransaksi: transaksi.IDTransaksi,
		IDObat:      req.IDObat,
		Jumlah:      req.Jumlah,
		HargaSatuan: hargaSatuan,
		Subtotal:    subtotal,
	}

	// Save detail
	if err := s.repo.CreateDetail(s.repo.GetDB(), &detailEntity); err != nil {
		return nil, fmt.Errorf("failed to add medicine to cart: %w", err)
	}

	// Update stock
	if err := s.repo.UpdateStokObat(s.repo.GetDB(), req.IDObat, req.Jumlah); err != nil {
		// Rollback detail
		s.repo.GetDB().Delete(&detailEntity)
		if err := s.repo.UpdateStatus(s.repo.GetDB(), transaksi.IDTransaksi, transaksi.Status); err != nil {
			return nil, fmt.Errorf("failed to rollback transaction: %w", err)
		}
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	// Reload transaction to get updated details
	if err := s.repo.GetDB().Preload("Details").First(&transaksi, transaksi.IDTransaksi).Error; err != nil {
		return nil, fmt.Errorf("failed to reload transaction: %w", err)
	}

	t := toTransactionResponse(transaksi)
	return &t, nil
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

func (s *transactionService) CancelTransaction(id uint, userID uint, userRole string) error {
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

	// Already cancelled
	if transaksi.Status == model.StatusTransaksiDibatalkan {
		tx.Rollback()
		return errors.New("transaction already cancelled")
	}

	isAdmin := userRole == "admin" || userRole == "owner"

	// Kasir (non-admin): hanya boleh batalkan transaksi miliknya sendiri
	if !isAdmin && transaksi.IDUser != userID {
		tx.Rollback()
		return errors.New("you do not have access to cancel this transaction")
	}

	// Kasir: hanya transaksi hari ini yang dapat dibatalkan.
	// Admin/owner dapat membatalkan transaksi hari apa pun (approval).
	if !isAdmin {
		now := time.Now()
		if transaksi.TglTransaksi.Year() != now.Year() ||
			transaksi.TglTransaksi.YearDay() != now.YearDay() {
			tx.Rollback()
			return errors.New(
				"transaction can only be cancelled on the same day",
			)
		}
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
