package returns

import (
	"errors"
	"fmt"
	"time"

	"medix-be/internal/transaction/model"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
)

type ReturnService interface {
	CreateReturn(userID uint, req dto.CreateReturnRequest) (*dto.CreateReturnResponse, error)
	ApproveReturn(adminID uint, returnID uint) (*dto.ApproveReturnResponse, error)
}

type returnService struct {
	repo repository.TransactionRepository
}

func NewReturnService(repo repository.TransactionRepository) ReturnService {
	return &returnService{repo: repo}
}

func (s *returnService) CreateReturn(userID uint, req dto.CreateReturnRequest) (*dto.CreateReturnResponse, error) {
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

		if items[i].KondisiLayak && ret.Status == "approved" {
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

func (s *returnService) ApproveReturn(adminID uint, returnID uint) (*dto.ApproveReturnResponse, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ret model.Return
	err := tx.Preload("Items").First(&ret, returnID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("retur tidak ditemukan")
	}

	if ret.Status != "pending" {
		tx.Rollback()
		return nil, errors.New("retur tidak dalam status pending")
	}

	user, err := s.repo.GetUser(tx, ret.DiajukanOleh)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}

	isAdmin := user.Role == "admin" || user.Role == "owner"
	if !isAdmin {
		tx.Rollback()
		return nil, errors.New("hanya admin yang dapat menyetujui retur")
	}

	if err := s.repo.UpdateReturnStatus(tx, returnID, "approved"); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Kembalikan stok untuk item yang layak
	var items []model.ReturnItem
	for _, item := range ret.Items {
		if item.KondisiLayak && item.StokKembali > 0 {
			if err := s.repo.RestoreStokObat(tx, item.IDObat, item.Jumlah); err != nil {
				tx.Rollback()
				return nil, err
			}
			item.StokKembali = item.Jumlah
			items = append(items, item)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.ApproveReturnResponse{
		IDReturn: ret.IDReturn,
		Status:   "approved",
		Message:  "retur berhasil disetujui dan stok dikembalikan",
	}, nil
}
