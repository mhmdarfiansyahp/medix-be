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

// US-15: Query agregasi grafik penjualan berdasarkan periode
func (r *reportRepository) GetSalesChart(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]dto.SalesChartData, error) {
	var results []dto.SalesChartData
	dateFormat := "YYYY-MM-DD"

	if groupBy == "monthly" {
		dateFormat = "YYYY-MM"
	}

	err := r.db.WithContext(ctx).
		Table("transaksi").
		Select("TO_CHAR(created_at, ?) as periode, COALESCE(SUM(total_harga), 0) as total_penjualan, COUNT(id_transaksi) as jumlah_transaksi", dateFormat).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Group("periode").
		Order("periode ASC").
		Scan(&results).Error

	return results, err
}

// US-16: Query Top Obat Terlaris
func (r *reportRepository) GetTopDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error) {
	var results []dto.DrugSalesStat

	err := r.db.WithContext(ctx).
		Table("detail_transaksi dt").
		Select("o.id_obat, o.nama_obat, COALESCE(SUM(dt.jumlah), 0) as total_terjual, COALESCE(SUM(dt.subtotal), 0) as total_omset").
		Joins("JOIN obat o ON o.id_obat = dt.obat_id").
		Joins("JOIN transaksi t ON t.id_transaksi = dt.transaksi_id").
		Where("t.created_at >= ? AND t.created_at <= ?", startDate, endDate).
		Group("o.id_obat, o.nama_obat").
		Order("total_terjual DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// US-16: Query Bottom Obat Terendah Terjual (termasuk obat yang 0 penjualan)
func (r *reportRepository) GetBottomDrugs(ctx context.Context, startDate, endDate time.Time, limit int) ([]dto.DrugSalesStat, error) {
	var results []dto.DrugSalesStat

	err := r.db.WithContext(ctx).
		Table("obat o").
		Select("o.id_obat, o.nama_obat, COALESCE(SUM(dt.jumlah), 0) as total_terjual, COALESCE(SUM(dt.subtotal), 0) as total_omset").
		Joins("LEFT JOIN detail_transaksi dt ON o.id_obat = dt.obat_id").
		Joins("LEFT JOIN transaksi t ON t.id_transaksi = dt.transaksi_id AND t.created_at >= ? AND t.created_at <= ?", startDate, endDate).
		Where("o.status = 1"). // Hanya mengambil obat yang aktif
		Group("o.id_obat, o.nama_obat").
		Order("total_terjual ASC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// US-17: Implementasi Query Data Ekspor Transaksi
func (r *reportRepository) GetExportTransactionData(ctx context.Context, startDate, endDate time.Time) ([]dto.ExportTransactionRow, error) {
	var results []dto.ExportTransactionRow

	err := r.db.WithContext(ctx).
		Table("detail_transaksi dt").
		Select("t.id_transaksi, t.created_at as tanggal, o.nama_obat, dt.harga_satuan, dt.jumlah, dt.subtotal").
		Joins("JOIN transaksi t ON t.id_transaksi = dt.transaksi_id").
		Joins("JOIN obat o ON o.id_obat = dt.obat_id").
		Where("t.created_at >= ? AND t.created_at <= ?", startDate, endDate).
		Order("t.created_at DESC").
		Scan(&results).Error

	return results, err
}
