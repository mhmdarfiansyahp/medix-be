package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"medix-be/internal/transaction/model"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
	"strconv"
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
	ProcessPayment(id uint, userID uint, req dto.PaymentRequest) (*dto.PaymentResponse, error)
	CreateReturn(userID uint, req dto.CreateReturnRequest) (*dto.CreateReturnResponse, error)
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
		MetodeBayar:  transaction.MetodeBayar,
		UangDiterima: transaction.UangDiterima,
		Kembalian:    transaction.Kembalian,
	}, nil
}

func (s *transactionService) ProcessPayment(id uint, userID uint, req dto.PaymentRequest) (*dto.PaymentResponse, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Cari transaksi
	var transaksi model.Transaksi
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
	if transaksi.Status != model.StatusTransaksiSelesai && transaksi.Status != model.StatusTransaksiProsesBayar {
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
	if err := s.repo.UpdateStatus(tx, id, model.StatusTransaksiProsesBayar); err != nil {
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

func (s *transactionService) CreateReturn(userID uint, req dto.CreateReturnRequest) (*dto.CreateReturnResponse, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Cari transaksi asli
	var transaksi model.Transaksi
	err := tx.
		Preload("Details").
		Where("id_transaksi = ?", req.IDTransaksi).
		First(&transaksi).
		Error

	if err != nil {
		tx.Rollback()
		return nil, errors.New("transaksi tidak ditemukan")
	}

	// Transaksi harus sudah selesai/final
	if transaksi.Status != model.StatusTransaksiSelesai && transaksi.Status != model.StatusTransaksiProsesBayar {
		tx.Rollback()
		return nil, errors.New("hanya transaksi yang sudah selesai yang dapat diretur")
	}

	// Validasi user: kasir hanya boleh retur transaksi miliknya sendiri, admin/owner boleh semua
	user, err := s.repo.GetUser(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}
	isAdmin := user.Role == "admin" || user.Role == "owner"
	if !isAdmin && transaksi.IDUser != userID {
		tx.Rollback()
		return nil, errors.New("anda hanya dapat meretur transaksi milik sendiri")
	}

	// Siapkan data retur dan validasi item
	ret := model.Return{
		IDTransaksi:  req.IDTransaksi,
		Alasan:       req.Alasan,
		DiajukanOleh: userID,
		Status:       "pending",
	}

	var totalNilaiRetur float64
	items := make([]model.ReturnItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		// Cari detail pembelian asli untuk validasi harga snapshot
		var detailAsli model.DetailPembelian
		err := tx.
			Where("id_transaksi = ? AND id_obat = ?", req.IDTransaksi, itemReq.IDObat).
			First(&detailAsli).
			Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("obat dengan id %d tidak ditemukan dalam transaksi", itemReq.IDObat)
		}

		if itemReq.Jumlah > detailAsli.Jumlah {
			tx.Rollback()
			return nil, fmt.Errorf("jumlah retur untuk obat %d melebihi jumlah pembelian", itemReq.IDObat)
		}

		nilaiRetur := float64(itemReq.Jumlah) * detailAsli.HargaSatuan
		totalNilaiRetur += nilaiRetur

		stokKembali := 0
		if itemReq.KondisiLayak {
			stokKembali = itemReq.Jumlah
		}

		item := model.ReturnItem{
			IDObat:       itemReq.IDObat,
			Jumlah:       itemReq.Jumlah,
			AlasanItem:   itemReq.AlasanItem,
			KondisiLayak: itemReq.KondisiLayak,
			StokKembali:  stokKembali,
		}
		items = append(items, item)
	}

	// Tentukan status otomatis: jika di bawah threshold dan semua item layak, langsung approved
	semuaLayak := true
	for _, it := range items {
		if !it.KondisiLayak {
			semuaLayak = false
			break
		}
	}

	if totalNilaiRetur <= 100000.0 && semuaLayak {
		ret.Status = "approved"
	}

	// Simpan header retur
	if err := s.repo.CreateReturn(tx, &ret); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create return: %w", err)
	}

	// Simpan item retur dan kembalikan stok jika layak
	for i := range items {
		items[i].IDReturn = ret.IDReturn
		if err := s.repo.CreateReturnItem(tx, &items[i]); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create return item: %w", err)
		}

		if items[i].KondisiLayak {
			if err := s.repo.RestoreStokObat(tx, items[i].IDObat, items[i].Jumlah); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to restore stock for medicine %d: %w", items[i].IDObat, err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Siapkan response
	respItems := make([]dto.ReturnItemResponse, 0, len(items))
	for _, it := range items {
		respItems = append(respItems, dto.ReturnItemResponse{
			IDItem:       it.IDItem,
			IDObat:       it.IDObat,
			Jumlah:       it.Jumlah,
			AlasanItem:   it.AlasanItem,
			KondisiLayak: it.KondisiLayak,
			StokKembali:  it.StokKembali,
		})
	}

	return &dto.CreateReturnResponse{
		IDReturn:     ret.IDReturn,
		IDTransaksi:  ret.IDTransaksi,
		Alasan:       ret.Alasan,
		TanggalRetur: ret.TanggalRetur.Format(time.RFC3339),
		DiajukanOleh: ret.DiajukanOleh,
		Status:       ret.Status,
		Items:        respItems,
	}, nil
}
