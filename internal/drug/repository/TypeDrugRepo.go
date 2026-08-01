package repository

import (
	"medix-be/internal/drug/model/entities"
	"gorm.io/gorm"
)

type TypeDrugRepository interface {
	Create(typeDrug *model.TypeDrug) error
	FindAll() ([]model.TypeDrug, error)
	FindByID(id uint) (*model.TypeDrug, error)
	Update(typeDrug *model.TypeDrug) error
	Delete(id uint) error
}

type typeDrugRepository struct {
	db *gorm.DB
}

func NewTypeDrugRepository(db *gorm.DB) TypeDrugRepository {
	return &typeDrugRepository{db: db}
}

func (r *typeDrugRepository) Create(typeDrug *model.TypeDrug) error {
	return r.db.Create(typeDrug).Error
}

func (r *typeDrugRepository) FindAll() ([]model.TypeDrug, error) {
	var typeDrugs []model.TypeDrug
	err := r.db.Find(&typeDrugs).Error
	return typeDrugs, err
}

func (r *typeDrugRepository) FindByID(id uint) (*model.TypeDrug, error) {
	var typeDrug model.TypeDrug
	err := r.db.First(&typeDrug, "id_jenis = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &typeDrug, nil
}

func (r *typeDrugRepository) Update(typeDrug *model.TypeDrug) error {
	return r.db.Save(typeDrug).Error
}

func (r *typeDrugRepository) Delete(id uint) error {
	return r.db.Delete(&model.TypeDrug{}, "id_jenis = ?", id).Error
}