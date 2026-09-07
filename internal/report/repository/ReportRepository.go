package repository

import (
	"context"
	"medix-be/internal/report/model/dto"
	"time"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetSalesChart(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]dto.SalesChartData, error)
	GetTopDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error)
	GetBottomDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error)
	GetExportTransactionData(ctx context.Context, startDate, endDate time.Time) ([]dto.ExportTransactionRow, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// US-15: Query sales chart aggregation by period
func (r *reportRepository) GetSalesChart(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]dto.SalesChartData, error) {
	var results []dto.SalesChartData
	dateFormat := "YYYY-MM-DD"

	switch groupBy {
	case "weekly":
		dateFormat = "YYYY-IW"
	case "monthly":
		dateFormat = "YYYY-MM"
	}

	err := r.db.WithContext(ctx).
		Table("transaksi").
		Select("TO_CHAR(created_at, ?) as periode, COALESCE(SUM(total_harga), 0) as total_penjualan, COUNT(id_transaksi) as jumlah_transaksi", dateFormat).
		Where("created_at >= ? AND created_at <= ? AND status = 1", startDate, endDate).
		Group("periode").
		Order("periode ASC").
		Scan(&results).Error

	return results, err
}

// US-16: Query Top Best-Selling Medicines
func (r *reportRepository) GetTopDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error) {
	var results []dto.DrugSalesStat

	err := r.db.WithContext(ctx).
		Table("detail_pembelian dt").
		Select("o.id_obat, o.nama_obat, COALESCE(SUM(dt.jumlah), 0) as total_terjual, COALESCE(SUM(dt.subtotal), 0) as total_omset").
		Joins("JOIN obat o ON o.id_obat = dt.id_obat").
		Joins("JOIN transaksi t ON t.id_transaksi = dt.id_transaksi").
		Where("t.created_at >= ? AND t.created_at <= ? AND t.status = 1", startDate, endDate).
		Group("o.id_obat, o.nama_obat").
		Order("total_terjual DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// US-16: Query Bottom Worst-Selling Medicines (including 0-sales medicines)
func (r *reportRepository) GetBottomDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error) {
	var results []dto.DrugSalesStat

	err := r.db.WithContext(ctx).
		Table("obat o").
		Select("o.id_obat, o.nama_obat, COALESCE(SUM(dt.jumlah), 0) as total_terjual, COALESCE(SUM(dt.subtotal), 0) as total_omset").
		Joins("LEFT JOIN detail_pembelian dt ON o.id_obat = dt.id_obat").
		Joins("LEFT JOIN transaksi t ON t.id_transaksi = dt.id_transaksi AND t.created_at >= ? AND t.created_at <= ? AND t.status = 1", startDate, endDate).
		Where("o.status = 1"). // Only fetch active medicines
		Group("o.id_obat, o.nama_obat").
		Order("total_terjual ASC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// US-17: Implement Export Transaction Data Query
func (r *reportRepository) GetExportTransactionData(ctx context.Context, startDate, endDate time.Time) ([]dto.ExportTransactionRow, error) {
	var results []dto.ExportTransactionRow

	err := r.db.WithContext(ctx).
		Table("detail_pembelian dt").
		Select("t.id_transaksi, t.created_at as tanggal, o.nama_obat, dt.harga_satuan, dt.jumlah, dt.subtotal").
		Joins("JOIN transaksi t ON t.id_transaksi = dt.id_transaksi").
		Joins("JOIN obat o ON o.id_obat = dt.id_obat").
		Where("t.created_at >= ? AND t.created_at <= ? AND t.status = 1", startDate, endDate).
		Order("t.created_at DESC").
		Scan(&results).Error

	return results, err
}
