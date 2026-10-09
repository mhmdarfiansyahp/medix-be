package report

import (
	"medix-be/internal/transaction/model"
	"medix-be/internal/transaction/model/dto"
	"medix-be/internal/transaction/repository"
)

type ReportService interface {
	GetTodayKasirReport(userID uint) (*dto.TodayKasirReportResponse, error)
}

type reportService struct {
	repo repository.TransactionRepository
}

func NewReportService(repo repository.TransactionRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetTodayKasirReport(userID uint) (*dto.TodayKasirReportResponse, error) {

	transactions, err := s.repo.FindTodayByUser(userID)
	if err != nil {
		return nil, err
	}

	totalPenjualan, totalTransaksi, err := s.repo.GetTodaySummary(userID)
	if err != nil {
		return nil, err
	}

	totalTunai, totalNonTunai, err := s.repo.GetTodayPaymentBreakdown(userID)
	if err != nil {
		return nil, err
	}

	returns, err := s.repo.FindTodayReturnsByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		result = append(result, toTransactionResponse(*transaction))
	}

	returnsRes := make([]dto.ReturnSummaryResponse, 0, len(returns))
	for _, r := range returns {
		returnsRes = append(returnsRes, dto.ReturnSummaryResponse{
			IDReturn:     r.IDReturn,
			IDTransaksi:  r.IDTransaksi,
			Alasan:       r.Alasan,
			TanggalRetur: r.TanggalRetur,
			Status:       r.Status,
		})
	}

	return &dto.TodayKasirReportResponse{
		Transactions: result,
		Returns:      returnsRes,
		Summary: dto.TransactionSummaryResponse{
			TotalTransaksi: totalTransaksi,
			TotalPenjualan: totalPenjualan,
			TotalTunai:     totalTunai,
			TotalNonTunai:  totalNonTunai,
		},
	}, nil
}

// Helper function (used by core service)
func toTransactionResponse(t model.Transaksi) dto.TransactionResponse {
	var detailsRes []dto.DetailItemResponse
	for _, d := range t.Details {
		detailsRes = append(detailsRes, dto.DetailItemResponse{
			IDDetail:    d.IDDetail,
			IDObat:      d.IDObat,
			Jumlah:      d.Jumlah,
			HargaSatuan: d.HargaSatuan,
			Subtotal:    d.Subtotal,
		})
	}

	return dto.TransactionResponse{
		IDTransaksi:  t.IDTransaksi,
		IDUser:       t.IDUser,
		TglTransaksi: t.TglTransaksi,
		TotalHarga:   t.TotalHarga,
		Status:       t.Status,
		Details:      detailsRes,
	}
}
