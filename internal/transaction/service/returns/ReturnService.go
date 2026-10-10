package returns

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"medix-be/internal/transaction/model/entities"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
)

type ReturnService interface {
	CreateReturn(userID uint, req dto.CreateReturnRequest) (*dto.CreateReturnResponse, error)
	ApproveReturn(adminID uint, returnID uint) (*dto.ApproveReturnResponse, error)
	GetReturnByID(returnID uint) (*dto.CreateReturnResponse, error)
	GetAllReturns() ([]*dto.CreateReturnResponse, error)
	RejectReturn(adminID uint, returnID uint, alasan string) (*dto.RejectReturnResponse, error)
	SetApprovalThreshold(threshold float64) (*dto.ApprovalThresholdResponse, error)
	GetApprovalThreshold() (*dto.ApprovalThresholdResponse, error)
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
	var transaksi entities.Transaksi
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
	if transaksi.Status != entities.StatusTransaksiSelesai && transaksi.Status != entities.StatusTransaksiProsesBayar {
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
	ret := entities.Return{
		IDTransaksi:  req.IDTransaksi,
		Alasan:       req.Alasan,
		DiajukanOleh: userID,
		Status:       "pending",
	}

	var totalNilaiRetur float64
	items := make([]entities.ReturnItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		// Cari detail pembelian asli untuk validasi harga snapshot
		var detailAsli entities.DetailPembelian
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

			item := entities.ReturnItem{
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

	if totalNilaiRetur <= getThreshold(s.repo) && semuaLayak {
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

	var ret entities.Return
	err := tx.Preload("Items").First(&ret, returnID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("retur tidak ditemukan")
	}

	if ret.Status != "pending" {
		tx.Rollback()
		return nil, errors.New("retur tidak dalam status pending")
	}

	if err := s.repo.UpdateReturnStatus(tx, returnID, "approved"); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Kembalikan stok untuk item yang layak
	for _, item := range ret.Items {
		if item.KondisiLayak && item.Jumlah > 0 {
			if err := s.repo.RestoreStokObat(tx, item.IDObat, item.Jumlah); err != nil {
				tx.Rollback()
				return nil, err
			}
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

func (s *returnService) GetReturnByID(returnID uint) (*dto.CreateReturnResponse, error) {
	ret, err := s.repo.GetReturnByID(returnID)
	if err != nil {
		return nil, errors.New("retur tidak ditemukan")
	}

	return s.buildReturnResponse(ret), nil
}

func (s *returnService) RejectReturn(adminID uint, returnID uint, alasan string) (*dto.RejectReturnResponse, error) {
	tx := s.repo.GetDB().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var ret entities.Return
	err := tx.First(&ret, returnID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("retur tidak ditemukan")
	}

	if ret.Status != "pending" {
		tx.Rollback()
		return nil, errors.New("retur tidak dalam status pending")
	}

	if err := s.repo.UpdateReturnStatusWithReason(tx, returnID, "rejected", alasan); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.RejectReturnResponse{
		IDReturn: returnID,
		Status:   "rejected",
		Alasan:   alasan,
		Message:  "retur berhasil ditolak",
	}, nil
}

func (s *returnService) SetApprovalThreshold(threshold float64) (*dto.ApprovalThresholdResponse, error) {
	if threshold < 0 {
		return nil, errors.New("batas approval tidak boleh negatif")
	}

	if err := s.repo.SetSetting("retur_approval_threshold", strconv.FormatFloat(threshold, 'f', -1, 64)); err != nil {
		return nil, err
	}

	return &dto.ApprovalThresholdResponse{
		Threshold: threshold,
		Message:   "Batas approval berhasil diatur",
	}, nil
}

func (s *returnService) GetApprovalThreshold() (*dto.ApprovalThresholdResponse, error) {
	threshold := getThreshold(s.repo)
	return &dto.ApprovalThresholdResponse{
		Threshold: threshold,
		Message:   "Batas approval retur",
	}, nil
}

func getThreshold(repo repository.TransactionRepository) float64 {
	value, err := repo.GetSetting("retur_approval_threshold")
	if err != nil {
		return 100000
	}
	threshold, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 100000
	}
	return threshold
}

func (s *returnService) GetAllReturns() ([]*dto.CreateReturnResponse, error) {
	returns, err := s.repo.FindAllReturns()
	if err != nil {
		return nil, err
	}

	resp := make([]*dto.CreateReturnResponse, 0, len(returns))
	for i := range returns {
		resp = append(resp, s.buildReturnResponse(&returns[i]))
	}

	return resp, nil
}

// buildReturnResponse maps a Return entity into its DTO, enriching items with
// the medicine name and purchase price snapshot from the original transaction.
// ponytail: N+1 lookups per return; batch-load if the returns list grows large.
func (s *returnService) buildReturnResponse(ret *entities.Return) *dto.CreateReturnResponse {
	db := s.repo.GetDB()

	var transaksi *entities.Transaksi
	if trx, err := s.repo.FindByID(ret.IDTransaksi); err == nil {
		transaksi = trx
	}

	respItems := make([]dto.ReturnItemResponse, 0, len(ret.Items))
	var total float64
	for _, it := range ret.Items {
		var harga float64
		if transaksi != nil {
			for _, d := range transaksi.Details {
				if d.IDObat == it.IDObat {
					harga = d.HargaSatuan
					break
				}
			}
		}

		nama := ""
		if obat, err := s.repo.FindObatByID(db, it.IDObat); err == nil {
			nama = obat.NamaObat
		}

		subtotal := harga * float64(it.Jumlah)
		total += subtotal

		respItems = append(respItems, dto.ReturnItemResponse{
			IDItem:       it.IDItem,
			IDObat:       it.IDObat,
			NamaObat:     nama,
			Jumlah:       it.Jumlah,
			HargaSatuan:  harga,
			Subtotal:     subtotal,
			AlasanItem:   it.AlasanItem,
			KondisiLayak: it.KondisiLayak,
			StokKembali:  it.StokKembali,
		})
	}

	return &dto.CreateReturnResponse{
		IDReturn:        ret.IDReturn,
		IDTransaksi:     ret.IDTransaksi,
		Alasan:          ret.Alasan,
		AlasanPenolakan: ret.AlasanPenolakan,
		TanggalRetur:    ret.TanggalRetur.Format(time.RFC3339),
		DiajukanOleh:    ret.DiajukanOleh,
		Status:          ret.Status,
		TotalNilai:      total,
		Items:           respItems,
	}
}
