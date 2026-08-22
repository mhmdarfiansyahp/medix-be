package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"medix-be/internal/report/model/dto"
	"medix-be/internal/report/repository"

	"github.com/xuri/excelize/v2"
)

type ReportService interface {
	GetSalesSummary(ctx context.Context, params dto.ReportFilterParams) (*dto.SalesSummaryResponse, error)
	GetDrugRanking(ctx context.Context, params dto.ReportFilterParams) (*dto.DrugRankingResponse, error)
	ExportToExcel(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetSalesSummary(ctx context.Context, params dto.ReportFilterParams) (*dto.SalesSummaryResponse, error) {
	start, end := parseDates(params.StartDate, params.EndDate)
	chartData, err := s.repo.GetSalesChart(ctx, start, end, params.GroupBy)
	if err != nil {
		return nil, err
	}

	var totalPenjualan float64
	var totalTransaksi int64
	for _, item := range chartData {
		totalPenjualan += item.TotalPenjualan
		totalTransaksi += item.JumlahTransaksi
	}

	return &dto.SalesSummaryResponse{
		TotalPenjualan: totalPenjualan,
		TotalTransaksi: totalTransaksi,
		ChartData:      chartData,
	}, nil
}

func (s *reportService) GetDrugRanking(ctx context.Context, params dto.ReportFilterParams) (*dto.DrugRankingResponse, error) {
	start, end := parseDates(params.StartDate, params.EndDate)

	top, err := s.repo.GetTopDrugs(ctx, start, end, 10)
	if err != nil {
		return nil, err
	}

	bottom, err := s.repo.GetBottomDrugs(ctx, start, end, 10)
	if err != nil {
		return nil, err
	}

	return &dto.DrugRankingResponse{
		TopMedicines:    top,
		BottomMedicines: bottom,
	}, nil
}

// US-17: Export Laporan ke format Excel (.xlsx)
func (s *reportService) ExportToExcel(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Laporan Penjualan"
	f.SetSheetName("Sheet1", sheet)

	// Header Table
	f.SetCellValue(sheet, "A1", "No")
	f.SetCellValue(sheet, "B1", "Nama Obat")
	f.SetCellValue(sheet, "C1", "Total Terjual")
	f.SetCellValue(sheet, "D1", "Total Omset (Rp)")

	start, end := parseDates(params.StartDate, params.EndDate)
	top, _ := s.repo.GetTopDrugs(ctx, start, end, 100) // Ambil data penjualan

	for i, row := range top {
		cellNum := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", cellNum), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", cellNum), row.NamaObat)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", cellNum), row.TotalTerjual)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", cellNum), row.TotalOmset)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

// Helper Parse Date Range Default 30 hari terakhir
func parseDates(startDate, endDate string) (time.Time, time.Time) {
	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)

	if err1 != nil {
		start = time.Now().AddDate(0, -1, 0) // Default 1 bulan lalu
	}
	if err2 != nil {
		end = time.Now()
	}

	// Atur jam akhir hari (23:59:59) untuk EndDate
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	return start, end
}
