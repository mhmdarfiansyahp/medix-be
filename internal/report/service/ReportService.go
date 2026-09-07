package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"medix-be/internal/report/model/dto"
	"medix-be/internal/report/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportService interface {
	GetSalesSummary(ctx context.Context, params dto.ReportFilterParams) (*dto.SalesSummaryResponse, error)
	GetDrugRanking(ctx context.Context, params dto.ReportFilterParams) (*dto.DrugRankingResponse, error)
	ExportToExcel(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error)
	ExportToPDF(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error)
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

func (s *reportService) ExportToExcel(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Sales Report"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "No")
	f.SetCellValue(sheet, "B1", "Date")
	f.SetCellValue(sheet, "C1", "Medicine")
	f.SetCellValue(sheet, "D1", "Unit Price (Rp)")
	f.SetCellValue(sheet, "E1", "Qty")
	f.SetCellValue(sheet, "F1", "Subtotal (Rp)")

	start, end := parseDates(params.StartDate, params.EndDate)
	rows, err := s.repo.GetExportTransactionData(ctx, start, end)
	if err != nil {
		return nil, err
	}

	for i, row := range rows {
		cellNum := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", cellNum), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", cellNum), row.Tanggal.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", cellNum), row.NamaObat)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", cellNum), row.HargaSatuan)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", cellNum), row.Jumlah)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", cellNum), row.Subtotal)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

func (s *reportService) ExportToPDF(ctx context.Context, params dto.ReportFilterParams) (*bytes.Buffer, error) {
	start, end := parseDates(params.StartDate, params.EndDate)
	rows, err := s.repo.GetExportTransactionData(ctx, start, end)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 14)
	pdf.Cell(0, 10, "Sales Report")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Period: %s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")))
	pdf.Ln(8)

	colW := []float64{10, 30, 60, 25, 15, 30}
	headers := []string{"No", "Date", "Medicine", "Unit Price", "Qty", "Subtotal"}

	pdf.SetFont("Helvetica", "B", 8)
	for i, h := range headers {
		pdf.CellFormat(colW[i], 6, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 7)
	totalQty := 0
	totalAmount := 0.0

	for i, row := range rows {
		pdf.CellFormat(colW[0], 6, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW[1], 6, row.Tanggal.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW[2], 6, truncate(row.NamaObat, 30), "1", 0, "L", false, 0, "")
		pdf.CellFormat(colW[3], 6, fmt.Sprintf("%.0f", row.HargaSatuan), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colW[4], 6, fmt.Sprintf("%d", row.Jumlah), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colW[5], 6, fmt.Sprintf("%.0f", row.Subtotal), "1", 0, "R", false, 0, "")
		pdf.Ln(6)

		totalQty += row.Jumlah
		totalAmount += row.Subtotal
	}

	pdf.SetFont("Helvetica", "B", 8)
	pdf.CellFormat(colW[0]+colW[1]+colW[2], 6, "TOTAL", "1", 0, "R", false, 0, "")
	pdf.CellFormat(colW[3], 6, fmt.Sprintf("%d", totalQty), "1", 0, "C", false, 0, "")
	pdf.CellFormat(colW[4]+colW[5], 6, fmt.Sprintf("%.0f", totalAmount), "1", 0, "R", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "I", 7)
	pdf.SetTextColor(128, 128, 128)
	pdf.Cell(0, 5, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func parseDates(startDate, endDate string) (time.Time, time.Time) {
	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)

	if err1 != nil {
		start = time.Now().AddDate(0, -1, 0)
	}
	if err2 != nil {
		end = time.Now()
	}

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	return start, end
}
