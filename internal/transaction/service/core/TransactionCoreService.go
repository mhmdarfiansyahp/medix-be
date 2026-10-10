package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"medix-be/internal/transaction/model/entities"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
)

type CoreService interface {
	CreateTransaction(userID uint, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	AddToCart(userID uint, req dto.AddToCartRequest) (*dto.TransactionResponse, error)
	GetAllTransactions() ([]dto.TransactionResponse, error)
	GetTransactionByID(id uint) (*dto.TransactionResponse, error)
	CancelTransaction(id uint, userID uint, userRole string) error
	ProcessPayment(id uint, userID uint, req dto.PaymentRequest) (*dto.PaymentResponse, error)
	GetTodayTransactions(userID uint) (*dto.TodayTransactionResponse, error)
}

type transactionCoreService struct {
	repo repository.TransactionRepository
}

func NewCoreService(repo repository.TransactionRepository) CoreService {
	return &transactionCoreService{repo: repo}
}

func (s *transactionCoreService) CreateTransaction(userID uint, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	var transaksi entities.Transaksi
	var detailsEntities []entities.DetailPembelian
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
				entities.DetailPembelian{
					IDObat:      item.IDObat,
					Jumlah:      item.Jumlah,
					HargaSatuan: hargaSatuan,
					Subtotal:    subtotal,
				},
			)
		}

		// Buat transaksi utama
		transaksi = entities.Transaksi{
			IDUser:     userID,
			TotalHarga: totalHarga,
			Status:     entities.StatusTransaksiSelesai,
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

func (s *transactionCoreService) AddToCart(userID uint, req dto.AddToCartRequest) (*dto.TransactionResponse, error) {
	var transaksi entities.Transaksi

	// Find today's transaction for this user (if any)
	var existingTransaksi entities.Transaksi
	error := s.repo.GetDB().
		Preload("Details").
		Where("id_user = ? AND status = ?", userID, entities.StatusTransaksiSelesai).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		First(&existingTransaksi).
		Error

	if error == nil && !errors.Is(error, gorm.ErrRecordNotFound) {
			// Found existing transaction
			transaksi = existingTransaksi
		} else {
			// Create new transaction
			transaksi = entities.Transaksi{
				IDUser:  userID,
				Status:  entities.StatusTransaksiSelesai,
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
	if err := s.repo.UpdateStatus(s.repo.GetDB(), transaksi.IDTransaksi, entities.StatusTransaksiSelesai); err != nil {
		return nil, fmt.Errorf("failed to update transaction total: %w", err)
	}

	// Create detail entity
	detailEntity := entities.DetailPembelian{
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

func (s *transactionCoreService) GetAllTransactions() ([]dto.TransactionResponse, error) {
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

func (s *transactionCoreService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	transaksi, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	res := toTransactionResponse(*transaksi)
	return &res, nil
}

func (s *transactionCoreService) CancelTransaction(id uint, userID uint, userRole string) error {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var transaksi entities.Transaksi

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
	if transaksi.Status == entities.StatusTransaksiDibatalkan {
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
	if err := s.repo.UpdateStatus(tx, id, entities.StatusTransaksiDibatalkan); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *transactionCoreService) ProcessPayment(id uint, userID uint, req dto.PaymentRequest) (*dto.PaymentResponse, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Cari transaksi
	var transaksi entities.Transaksi
	err := tx.
		Preload("Details").
		Where("id_transaksi = ?", id).
		First(&transaksi).
		Error

	if err != nil {
		tx.Rollback()
		return nil, errors.New("transaksi tidak ditemukan")
	}

	// Validasi status transaksi (harus belum dibayar/final)
	if transaksi.Status != entities.StatusTransaksiSelesai && transaksi.Status != entities.StatusTransaksiProsesBayar {
		tx.Rollback()
		return nil, errors.New("transaksi tidak dapat dibayar, status saat ini: " + strconv.Itoa(transaksi.Status))
	}

	// Validasi metode pembayaran
	if req.MetodeBayar == "tunai" {
		if req.UangDiterima <= 0 {
			tx.Rollback()
			return nil, errors.New("uang diterima harus lebih besar dari 0 untuk pembayaran tunai")
		}
		if req.UangDiterima < transaksi.TotalHarga {
			tx.Rollback()
			return nil, errors.New("uang diterima tidak mencukupi untuk total transaksi")
		}
		transaksi.UangDiterima = req.UangDiterima
		transaksi.Kembalian = req.UangDiterima - transaksi.TotalHarga
	} else {
		// Pembayaran non-tunai tidak memerlukan input uang diterima
		transaksi.UangDiterima = 0
		transaksi.Kembalian = 0
	}
	transaksi.MetodeBayar = req.MetodeBayar
	
	// Update status transaksi menjadi final (status 2)
	if err := s.repo.UpdateStatus(tx, id, entities.StatusTransaksiProsesBayar); err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload transaksi untuk mendapatkan data terbaru
	if err := tx.Preload("Details").First(&transaksi, id).Error; err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		IDTransaksi:  transaksi.IDTransaksi,
		MetodeBayar:    transaksi.MetodeBayar,
		UangDiterima:   transaksi.UangDiterima,
		Kembalian:      transaksi.Kembalian,
		Status:       transaksi.Status,
		TotalHarga:   transaksi.TotalHarga,
	}, nil
}

func (s *transactionCoreService) GetTodayTransactions(userID uint) (*dto.TodayTransactionResponse, error) {

	transactions, err := s.repo.FindTodayByUser(userID)
	if err != nil {
		return nil, err
	}

	totalPenjualan, totalTransaksi, err :=
		s.repo.GetTodaySummary(userID)

	if err != nil {
		return nil, err
	}

	// Calculate payment breakdown
	totalTunai, totalNonTunai, err := s.repo.GetTodayPaymentBreakdown(userID)
	if err != nil {
		return nil, err
	}

	// Get today's returns
	returns, err := s.repo.FindTodayReturnsByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		result = append(result, toTransactionResponse(*transaction))
	}

	returnsRes := make([]dto.ReturnSummaryResponse, 0, len(returns))
	for _, r := range returns {
		returnsRes = append(returnsRes, dto.ReturnSummaryResponse{
			IDReturn:     r.IDReturn,
			IDTransaksi:  r.IDTransaksi,
			Alasan:       r.Alasan,
			TanggalRetur: r.TanggalRetur,
			Status:       r.Status,
		})
	}

	return &dto.TodayTransactionResponse{
		Transactions: result,
		Summary: dto.TransactionSummaryResponse{
			TotalTransaksi: totalTransaksi,
			TotalPenjualan: totalPenjualan,
			TotalTunai:     totalTunai,
			TotalNonTunai:  totalNonTunai,
		},
		Returns: returnsRes,
	}, nil
}

func toTransactionResponse(t entities.Transaksi) dto.TransactionResponse {
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
