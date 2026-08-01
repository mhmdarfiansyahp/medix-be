package service

import (
	"errors"
	"medix-be/internal/drug/model/dto"
	"medix-be/internal/drug/model/entities"
	"medix-be/internal/drug/repository"
)

type TypeDrugService interface {
	CreateTypeDrug(req dto.CreateTypeDrugRequest) (*dto.TypeDrugResponse, error)
	GetAllTypeDrugs() ([]dto.TypeDrugResponse, error)
	GetTypeDrugByID(id uint) (*dto.TypeDrugResponse, error)
	UpdateTypeDrug(id uint, req dto.UpdateTypeDrugRequest) (*dto.TypeDrugResponse, error)
	DeleteTypeDrug(id uint) error
}

type typeDrugService struct {
	repo repository.TypeDrugRepository
}

func NewTypeDrugService(repo repository.TypeDrugRepository) TypeDrugService {
	return &typeDrugService{repo: repo}
}

func (s *typeDrugService) CreateTypeDrug(req dto.CreateTypeDrugRequest) (*dto.TypeDrugResponse, error) {
	typeDrug := model.TypeDrug{
		NamaJenis: req.NamaJenis,
	}

	if err := s.repo.Create(&typeDrug); err != nil {
		return nil, err
	}

	res := toTypeDrugResponse(typeDrug)
	return &res, nil
}

func (s *typeDrugService) GetAllTypeDrugs() ([]dto.TypeDrugResponse, error) {
	typeDrugs, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.TypeDrugResponse
	for _, td := range typeDrugs {
		responses = append(responses, toTypeDrugResponse(td))
	}
	return responses, nil
}

func (s *typeDrugService) GetTypeDrugByID(id uint) (*dto.TypeDrugResponse, error) {
	typeDrug, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("jenis obat tidak ditemukan")
	}

	res := toTypeDrugResponse(*typeDrug)
	return &res, nil
}

func (s *typeDrugService) UpdateTypeDrug(id uint, req dto.UpdateTypeDrugRequest) (*dto.TypeDrugResponse, error) {
	typeDrug, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("jenis obat tidak ditemukan")
	}

	if req.NamaJenis != "" {
		typeDrug.NamaJenis = req.NamaJenis
	}

	if err := s.repo.Update(typeDrug); err != nil {
		return nil, err
	}

	res := toTypeDrugResponse(*typeDrug)
	return &res, nil
}

func (s *typeDrugService) DeleteTypeDrug(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("jenis obat tidak ditemukan")
	}

	return s.repo.Delete(id)
}

// Helper untuk format mapping response
func toTypeDrugResponse(td model.TypeDrug) dto.TypeDrugResponse {
	return dto.TypeDrugResponse{
		IDJenis:   td.IDJenis,
		NamaJenis: td.NamaJenis,
	}
}