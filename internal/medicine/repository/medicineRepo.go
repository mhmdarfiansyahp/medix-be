package repository

import (
	"medix-be/internal/medicine/model/entities"
	"gorm.io/gorm"
)

type MedicineRepository interface {
	Create(medicine *model.Obat) error
	FindAll() ([]model.Obat, error)
	FindByID(id uint) (*model.Obat, error)
	FindByBarcode(barcode string) (*model.Obat, error)
	Update(medicine *model.Obat) error
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

func (r *medicineRepository) FindAll() ([]model.Obat, error) {
	var medicines []model.Obat
	err := r.db.Preload("JenisObat").Find(&medicines).Error
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

func (r *medicineRepository) Delete(id uint) error {
	return r.db.Delete(&model.Obat{}, "id_obat = ?", id).Error
}