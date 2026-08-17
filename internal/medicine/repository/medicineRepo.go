package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"medix-be/internal/medicine/model/dto"
	"medix-be/internal/medicine/model/entities"
)

type MedicineRepository interface {
	Create(medicine *model.Obat) error
	FindAll(ctx context.Context, params dto.MedicineFilterParams) ([]model.Obat, error)
	FindByID(id uint) (*model.Obat, error)
	FindByBarcode(barcode string) (*model.Obat, error)
	Update(medicine *model.Obat) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	Delete(id uint) error
}

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &medicineRepository{db: db}
}

func (r *medicineRepository) Create(medicine *model.Obat) error {
	return r.db.Create(medicine).Error
}

func (r *medicineRepository) FindAll(ctx context.Context, params dto.MedicineFilterParams) ([]model.Obat, error) {
	var medicines []model.Obat
	query := r.db.WithContext(ctx).Preload("JenisObat")

	if params.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where("LOWER(nama_obat) LIKE LOWER(?)", searchTerm)
	}

	if params.JenisObatID > 0 {
		query = query.Where("jenis_obat_id = ?", params.JenisObatID)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.StatusStok != "" {
		switch params.StatusStok {
		case "habis":
			query = query.Where("stok = 0")
		case "menipis":
			query = query.Where("stok <= stok_minimum AND stok > 0")
		case "tersedia":
			query = query.Where("stok > 0")
		}
	}

	if params.Barcode != "" {
		query = query.Where("barcode = ?", params.Barcode)
	}

	err := query.Find(&medicines).Error
	return medicines, err
}

func (r *medicineRepository) FindByID(id uint) (*model.Obat, error) {
	var medicine model.Obat
	err := r.db.Preload("JenisObat").First(&medicine, "id_obat = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *medicineRepository) FindByBarcode(barcode string) (*model.Obat, error) {
	var medicine model.Obat
	err := r.db.Preload("JenisObat").First(&medicine, "barcode = ?", barcode).Error
	if err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (r *medicineRepository) Update(medicine *model.Obat) error {
	return r.db.Save(medicine).Error
}

func (r *medicineRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).Model(&model.Obat{}).Where("id_obat = ?", id).Update("status", status).Error
}

func (r *medicineRepository) Delete(id uint) error {
	return r.db.Delete(&model.Obat{}, "id_obat = ?", id).Error
}
