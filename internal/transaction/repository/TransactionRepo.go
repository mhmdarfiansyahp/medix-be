package repository

import (
	"errors"

	"gorm.io/gorm"

	medicine "medix-be/internal/medicine/model/entities"
	"medix-be/internal/transaction/model/entities"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaksi *entities.Transaksi) error
	CreateDetail(tx *gorm.DB, detail *entities.DetailPembelian) error

	// Obat
	FindObatByID(tx *gorm.DB, idObat uint) (*medicine.Obat, error)
	UpdateStokObat(tx *gorm.DB, idObat uint, jumlah int) error
	RestoreStokObat(tx *gorm.DB, idObat uint, jumlah int) error

	// Transaksi
	FindAll() ([]entities.Transaksi, error)
	FindByID(id uint) (*entities.Transaksi, error)
	FindTodayByUser(idUser uint) ([]*entities.Transaksi, error)
	GetTodaySummary(idUser uint) (float64, int, error)
	GetTodayPaymentBreakdown(idUser uint) (float64, float64, error)
	UpdateStatus(tx *gorm.DB, id uint, status int) error
	GetUser(tx *gorm.DB, id uint) (*entities.User, error)

	// Retur
	CreateReturn(tx *gorm.DB, ret *entities.Return) error
	CreateReturnItem(tx *gorm.DB, item *entities.ReturnItem) error
	GetReturnByID(id uint) (*entities.Return, error)
	FindTodayReturnsByUser(idUser uint) ([]entities.Return, error)
	UpdateReturnStatus(tx *gorm.DB, id uint, status string) error
	UpdateReturnStatusWithReason(tx *gorm.DB, id uint, status string, reason string) error

	// Pengaturan
	GetSetting(key string) (string, error)
	SetSetting(key string, value string) error

	GetDB() *gorm.DB
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *transactionRepository) Create(tx *gorm.DB, transaksi *entities.Transaksi) error {
	return tx.Create(transaksi).Error
}

func (r *transactionRepository) CreateDetail(tx *gorm.DB, detail *entities.DetailPembelian) error {
	return tx.Create(detail).Error
}

func (r *transactionRepository) FindObatByID(tx *gorm.DB, idObat uint) (*medicine.Obat, error) {

	var obat medicine.Obat

	err := tx.
		Where("id_obat = ?", idObat).
		First(&obat).
		Error

	if err != nil {
		return nil, err
	}

	return &obat, nil
}

func (r *transactionRepository) UpdateStokObat(tx *gorm.DB, idObat uint, jumlah int) error {
	result := tx.
		Table("obat").
		Where(
			"id_obat = ? AND stok >= ? AND status = 1",
			idObat,
			jumlah,
		).
		UpdateColumn(
			"stok",
			gorm.Expr("stok - ?", jumlah),
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New(
			"stok obat tidak mencukupi atau obat tidak aktif",
		)
	}

	return nil
}

func (r *transactionRepository) RestoreStokObat(tx *gorm.DB, idObat uint, jumlah int) error {

	result := tx.
		Table("obat").
		Where("id_obat = ?", idObat).
		UpdateColumn(
			"stok",
			gorm.Expr("stok + ?", jumlah),
		)

	return result.Error
}

func (r *transactionRepository) FindAll() ([]entities.Transaksi, error) {

	var list []entities.Transaksi

	err := r.db.
		Preload("Details").
		Find(&list).
		Error

	return list, err
}

func (r *transactionRepository) FindByID(id uint) (*entities.Transaksi, error) {

	var transaksi entities.Transaksi

	err := r.db.
		Preload("Details").
		First(&transaksi, id).
		Error

	if err != nil {
		return nil, err
	}

	return &transaksi, nil
}

func (r *transactionRepository) FindTodayByUser(idUser uint) ([]*entities.Transaksi, error) {

	var list []*entities.Transaksi

	err := r.db.
		Preload("Details").
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Order("tgl_transaksi DESC").
		Find(&list).
		Error

	return list, err
}

func (r *transactionRepository) GetTodaySummary(idUser uint) (float64, int, error) {

	var total float64
	var count int64

	err := r.db.
		Model(&entities.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", entities.StatusTransaksiSelesai).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&total).
		Error

	if err != nil {
		return 0, 0, err
	}

	err = r.db.
		Model(&entities.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", entities.StatusTransaksiSelesai).
		Count(&count).
		Error

	if err != nil {
		return 0, 0, err
	}

	return total, int(count), nil
}

func (r *transactionRepository) UpdateStatus(tx *gorm.DB, id uint, status int) error {

	return tx.
		Model(&entities.Transaksi{}).
		Where("id_transaksi = ?", id).
		Update("status", status).
		Error
}

func (r *transactionRepository) GetTodayPaymentBreakdown(idUser uint) (float64, float64, error) {
	var tunai, nonTunai float64
	err := r.db.
		Model(&entities.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", entities.StatusTransaksiSelesai).
		Where("metode_bayar = ?", "tunai").
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&tunai).
		Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.
		Model(&entities.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", entities.StatusTransaksiSelesai).
		Where("metode_bayar != ?", "tunai").
		Where("metode_bayar IS NOT NULL").
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&nonTunai).
		Error
	if err != nil {
		return 0, 0, err
	}

	return tunai, nonTunai, nil
}

func (r *transactionRepository) GetUser(tx *gorm.DB, id uint) (*entities.User, error) {
	var user entities.User
	err := tx.
		Where("id_user = ?", id).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *transactionRepository) FindTodayReturnsByUser(idUser uint) ([]entities.Return, error) {
	var returns []entities.Return
	err := r.db.
		Where("diajukan_oleh = ?", idUser).
		Where("DATE(tanggal_retur) = CURRENT_DATE").
		Preload("Items").
		Order("tanggal_retur DESC").
		Find(&returns).
		Error
	if err != nil {
		return nil, err
	}
	return returns, nil
}

func (r *transactionRepository) CreateReturn(tx *gorm.DB, ret *entities.Return) error {
	return tx.Create(ret).Error
}

func (r *transactionRepository) CreateReturnItem(tx *gorm.DB, item *entities.ReturnItem) error {
	return tx.Create(item).Error
}

func (r *transactionRepository) GetReturnByID(id uint) (*entities.Return, error) {
	var ret entities.Return
	err := r.db.
		Preload("Items").
		Where("id_return = ?", id).
		First(&ret).
		Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (r *transactionRepository) UpdateReturnStatus(tx *gorm.DB, id uint, status string) error {
	return tx.
		Model(&entities.Return{}).
		Where("id_return = ?", id).
		Update("status", status).
		Error
}

func (r *transactionRepository) UpdateReturnStatusWithReason(tx *gorm.DB, id uint, status string, reason string) error {
	return tx.
		Model(&entities.Return{}).
		Where("id_return = ?", id).
		Updates(map[string]interface{}{
			"status":           status,
			"alasan_penolakan": reason,
		}).
		Error
}

func (r *transactionRepository) GetSetting(key string) (string, error) {
	var s entities.Pengaturan
	err := r.db.
		Where("key = ?", key).
		First(&s).
		Error
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

func (r *transactionRepository) SetSetting(key string, value string) error {
	return r.db.
		Model(&entities.Pengaturan{}).
		Where("key = ?", key).
		Update("value", value).
		Error
}
