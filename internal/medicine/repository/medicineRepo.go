package repository

import (
	"context"
	"fmt"
	"medix-be/internal/medicine/model/dto"
	"medix-be/internal/medicine/model/entities"
	"time"

	"gorm.io/gorm"
)

type MedicineRepository interface {
	Create(medicine *model.Obat) error
	FindAll(ctx context.Context, params dto.MedicineFilterParams) ([]model.Obat, error)
	FindByID(id uint) (*model.Obat, error)
	FindByBarcode(barcode string) (*model.Obat, error)
	Update(medicine *model.Obat) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	Delete(id uint) error

	GetLowStock(ctx context.Context) ([]dto.LowStockResponse, error)
	GetExpiring(ctx context.Context, days int) ([]dto.ExpiringDrugResponse, error)
	GetNotificationSummary(ctx context.Context) (dto.StockNotificationSummary, error)
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

		query = query.Where(
			r.db.Where("LOWER(nama_obat) LIKE LOWER(?)", searchTerm).
				Or("LOWER(COALESCE(barcode, '')) LIKE LOWER(?)", searchTerm).
				Or("LOWER(COALESCE(merk_obat, '')) LIKE LOWER(?)", searchTerm),
		)
	}

	if params.JenisObatID > 0 {
		query = query.Where("jenis_obat_id = ?", params.JenisObatID)
	}

	if params.Status != "" {
		if params.Status == "1" || params.Status == "active" || params.Status == "true" {
			query = query.Where("status = 1")
		} else if params.Status == "0" || params.Status == "inactive" || params.Status == "false" {
			query = query.Where("status = 0")
		}
	}

	if params.StatusStok != "" {
		switch params.StatusStok {
		case "habis", "out":
			query = query.Where("stok = 0")
		case "menipis", "low":
			query = query.Where("stok <= stok_minimum AND stok > 0")
		case "tersedia", "available":
			query = query.Where("stok > 0")
		}
	}

	if params.Barcode != "" {
		query = query.Where("barcode = ?", params.Barcode)
	}

	if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
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

func (r *medicineRepository) GetLowStock(ctx context.Context) ([]dto.LowStockResponse, error) {
	var results []dto.LowStockResponse

	err := r.db.WithContext(ctx).
		Table("obat o").
		Select("o.id_obat, o.nama_obat, o.merk_obat, o.stok, o.stok_minimum, COALESCE(j.nama_jenis, '-') as nama_jenis").
		Joins("LEFT JOIN jenis_obat j ON j.id_jenis = o.jenis_obat_id").
		Where("o.stok <= o.stok_minimum AND o.status = 1").
		Order("o.stok ASC").
		Scan(&results).Error

	return results, err
}

// US-14: Mengambil daftar obat yang mendekati tanggal kadaluarsa (misal 30 hari ke depan)
func (r *medicineRepository) GetExpiring(ctx context.Context, days int) ([]dto.ExpiringDrugResponse, error) {
	var results []dto.ExpiringDrugResponse
	targetDate := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Table("obat o").
		Select("o.id_obat, o.nama_obat, o.merk_obat, o.tgl_kadaluarsa, o.stok, COALESCE(j.nama_jenis, '-') as nama_jenis, DATE_PART('day', o.tgl_kadaluarsa - NOW())::INT as sisa_hari").
		Joins("LEFT JOIN jenis_obat j ON j.id_jenis = o.jenis_obat_id").
		Where("o.tgl_kadaluarsa >= CURRENT_DATE AND o.tgl_kadaluarsa <= ? AND o.status = 1", targetDate).
		Order("o.tgl_kadaluarsa ASC").
		Scan(&results).Error

	return results, err
}

func (r *medicineRepository) GetNotificationSummary(ctx context.Context) (dto.StockNotificationSummary, error) {
	var summary dto.StockNotificationSummary
	targetDate := time.Now().AddDate(0, 0, 30)

	r.db.WithContext(ctx).
		Table("obat").
		Where("stok <= stok_minimum AND status = 1").
		Count(&summary.TotalLowStock)

	r.db.WithContext(ctx).
		Table("obat").
		Where("tgl_kadaluarsa >= CURRENT_DATE AND tgl_kadaluarsa <= ? AND status = 1", targetDate).
		Count(&summary.TotalExpiring)

	return summary, nil
}
