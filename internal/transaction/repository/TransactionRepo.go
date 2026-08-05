package repository

import (
	"medix-be/internal/transaction/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaksi *model.Transaksi) error
	CreateDetail(tx *gorm.DB, detail *model.DetailPembelian) error
	UpdateStokObat(tx *gorm.DB, idObat uint, jumlah int) error
	FindAll() ([]model.Transaksi, error)
	FindByID(id uint) (*model.Transaksi, error)
	GetDB() *gorm.DB
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
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

func (r *transactionRepository) UpdateStokObat(tx *gorm.DB, idObat uint, jumlah int) error {
	// Memotong stok obat secara aman
	return tx.Table("obat").Where("id_obat = ? AND stok >= ?", idObat, jumlah).
		UpdateColumn("stok", gorm.Expr("stok - ?", jumlah)).Error
}

func (r *transactionRepository) FindAll() ([]model.Transaksi, error) {
	var list []model.Transaksi
	err := r.db.Preload("Details").Find(&list).Error
	return list, err
}

func (r *transactionRepository) FindByID(id uint) (*model.Transaksi, error) {
	var transaksi model.Transaksi
	err := r.db.Preload("Details").First(&transaksi, id).Error
	return &transaksi, err
}