package service

import (
	"context"
	"errors"
	"medix-be/internal/medicine/model/dto"
	"medix-be/internal/medicine/model/entities"
	"medix-be/internal/medicine/repository"
	"time"
)

type MedicineService interface {
	CreateMedicine(req dto.CreateMedicineRequest) (*dto.MedicineResponse, error)
	GetAllMedicines(ctx context.Context, params dto.MedicineFilterParams) ([]dto.MedicineResponse, error)
	GetMedicineByID(id uint) (*dto.MedicineResponse, error)
	GetMedicineByBarcode(ctx context.Context, barcode string) (*dto.MedicineResponse, error)
	UpdateMedicine(id uint, req dto.UpdateMedicineRequest) (*dto.MedicineResponse, error)
	ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error
	DeleteMedicine(id uint) error
	GetLowStockDrugs(ctx context.Context) ([]dto.LowStockResponse, error)
	GetExpiringDrugs(ctx context.Context, days int) ([]dto.ExpiringDrugResponse, error)
	GetNotificationSummary(ctx context.Context) (*dto.StockNotificationSummary, error)
}

type medicineService struct {
	repo repository.MedicineRepository
}

func NewMedicineService(repo repository.MedicineRepository) MedicineService {
	return &medicineService{repo: repo}
}

func (s *medicineService) CreateMedicine(req dto.CreateMedicineRequest) (*dto.MedicineResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", req.TglKadaluarsa)
	if err != nil {
		return nil, errors.New("format tgl_kadaluarsa harus YYYY-MM-DD")
	}

	if req.Barcode != nil && *req.Barcode != "" {
		existing, _ := s.repo.FindByBarcode(*req.Barcode)
		if existing != nil {
			return nil, errors.New("barcode sudah digunakan oleh obat lain")
		}
	}

	stokMin := req.StokMinimum
	if stokMin == 0 {
		stokMin = 10
	}

	medicine := model.Obat{
		NamaObat:      req.NamaObat,
		MerkObat:      req.MerkObat,
		JenisObatID:   req.JenisObatID,
		Barcode:       req.Barcode,
		TglKadaluarsa: parsedDate,
		Harga:         req.Harga,
		Stok:          req.Stok,
		StokMinimum:   stokMin,
		Keterangan:    req.Keterangan,
		Status:        1,
		Gambar:        req.Gambar,
	}

	if err := s.repo.Create(&medicine); err != nil {
		return nil, err
	}

	res := toMedicineResponse(medicine)
	return &res, nil
}

func (s *medicineService) GetAllMedicines(ctx context.Context, params dto.MedicineFilterParams) ([]dto.MedicineResponse, error) {
	medicines, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	var responses []dto.MedicineResponse
	for _, m := range medicines {
		responses = append(responses, toMedicineResponse(m))
	}
	return responses, nil
}

func (s *medicineService) GetMedicineByID(id uint) (*dto.MedicineResponse, error) {
	medicine, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("obat tidak ditemukan")
	}

	res := toMedicineResponse(*medicine)
	return &res, nil
}

func (s *medicineService) GetMedicineByBarcode(ctx context.Context, barcode string) (*dto.MedicineResponse, error) {
	medicine, err := s.repo.FindByBarcode(barcode)
	if err != nil {
		return nil, errors.New("obat dengan barcode tersebut tidak ditemukan")
	}

	res := toMedicineResponse(*medicine)
	return &res, nil
}

func (s *medicineService) UpdateMedicine(id uint, req dto.UpdateMedicineRequest) (*dto.MedicineResponse, error) {
	medicine, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("obat tidak ditemukan")
	}

	if req.Barcode != nil && *req.Barcode != "" {
		existing, _ := s.repo.FindByBarcode(*req.Barcode)
		if existing != nil && existing.IDObat != id {
			return nil, errors.New("barcode sudah digunakan oleh obat lain")
		}
		medicine.Barcode = req.Barcode
	}

	if req.NamaObat != "" {
		medicine.NamaObat = req.NamaObat
	}
	if req.MerkObat != "" {
		medicine.MerkObat = req.MerkObat
	}
	if req.JenisObatID != 0 {
		medicine.JenisObatID = req.JenisObatID
	}
	if req.TglKadaluarsa != "" {
		parsedDate, err := time.Parse("2006-01-02", req.TglKadaluarsa)
		if err != nil {
			return nil, errors.New("format tgl_kadaluarsa harus YYYY-MM-DD")
		}
		medicine.TglKadaluarsa = parsedDate
	}
	if req.Harga >= 0 {
		medicine.Harga = req.Harga
	}
	if req.Stok >= 0 {
		medicine.Stok = req.Stok
	}
	if req.StokMinimum > 0 {
		medicine.StokMinimum = req.StokMinimum
	}
	if req.Keterangan != "" {
		medicine.Keterangan = req.Keterangan
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("status harus 0 atau 1")
		}

		medicine.Status = *req.Status
	}
	if req.Gambar != "" {
		medicine.Gambar = req.Gambar
	}

	if err := s.repo.Update(medicine); err != nil {
		return nil, err
	}

	res := toMedicineResponse(*medicine)
	return &res, nil
}

func (s *medicineService) ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("obat tidak ditemukan")
	}

	status := 0
	if isActive {
		status = 1
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *medicineService) DeleteMedicine(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("obat tidak ditemukan")
	}

	return s.repo.Delete(id)
}

func (s *medicineService) GetLowStockDrugs(ctx context.Context) ([]dto.LowStockResponse, error) {
	return s.repo.GetLowStock(ctx)
}

func (s *medicineService) GetExpiringDrugs(ctx context.Context, days int) ([]dto.ExpiringDrugResponse, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetExpiring(ctx, days)
}

func (s *medicineService) GetNotificationSummary(ctx context.Context) (*dto.StockNotificationSummary, error) {
	summary, err := s.repo.GetNotificationSummary(ctx)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// Helper Mapping Response
func toMedicineResponse(m model.Obat) dto.MedicineResponse {
	return dto.MedicineResponse{
		IDObat:        m.IDObat,
		NamaObat:      m.NamaObat,
		MerkObat:      m.MerkObat,
		JenisObatID:   m.JenisObatID,
		JenisObat:     m.JenisObat,
		Barcode:       m.Barcode,
		TglKadaluarsa: m.TglKadaluarsa.Format("2006-01-02"),
		Harga:         m.Harga,
		Stok:          m.Stok,
		StokMinimum:   m.StokMinimum,
		Keterangan:    m.Keterangan,
		Status:        m.Status,
		Gambar:        m.Gambar,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
