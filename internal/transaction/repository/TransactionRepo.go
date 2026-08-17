package repository

import (
	"errors"

	"gorm.io/gorm"

	medicine "medix-be/internal/medicine/model/entities"
	"medix-be/internal/transaction/model"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaksi *model.Transaksi) error
	CreateDetail(tx *gorm.DB, detail *model.DetailPembelian) error

	// Obat
	FindObatByID(tx *gorm.DB, idObat uint) (*medicine.Obat, error)
	UpdateStokObat(tx *gorm.DB, idObat uint, jumlah int) error
	RestoreStokObat(tx *gorm.DB, idObat uint, jumlah int) error

	// Transaksi
	FindAll() ([]model.Transaksi, error)
	FindByID(id uint) (*model.Transaksi, error)
	FindTodayByUser(idUser uint) ([]*model.Transaksi, error)
	GetTodaySummary(idUser uint) (float64, int, error)
	UpdateStatus(tx *gorm.DB, id uint, status int) error

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

func (r *transactionRepository) Create(tx *gorm.DB, transaksi *model.Transaksi) error {
	return tx.Create(transaksi).Error
}

func (r *transactionRepository) CreateDetail(tx *gorm.DB, detail *model.DetailPembelian) error {
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

func (r *transactionRepository) FindAll() ([]model.Transaksi, error) {

	var list []model.Transaksi

	err := r.db.
		Preload("Details").
		Find(&list).
		Error

	return list, err
}

func (r *transactionRepository) FindByID(id uint) (*model.Transaksi, error) {

	var transaksi model.Transaksi

	err := r.db.
		Preload("Details").
		First(&transaksi, id).
		Error

	if err != nil {
		return nil, err
	}

	return &transaksi, nil
}

func (r *transactionRepository) FindTodayByUser(idUser uint) ([]*model.Transaksi, error) {

	var list []*model.Transaksi

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
		Model(&model.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", model.StatusTransaksiSelesai).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&total).
		Error

	if err != nil {
		return 0, 0, err
	}

	err = r.db.
		Model(&model.Transaksi{}).
		Where("id_user = ?", idUser).
		Where("DATE(tgl_transaksi) = CURRENT_DATE").
		Where("status = ?", model.StatusTransaksiSelesai).
		Count(&count).
		Error

	if err != nil {
		return 0, 0, err
	}

	return total, int(count), nil
}

func (r *transactionRepository) UpdateStatus(tx *gorm.DB, id uint, status int) error {

	return tx.
		Model(&model.Transaksi{}).
		Where("id_transaksi = ?", id).
		Update("status", status).
		Error
}
