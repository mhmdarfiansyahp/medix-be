package dto

import "time"

type ExportTransactionRow struct {
	IDTransaksi uint      `json:"id_transaksi"`
	Tanggal     time.Time `json:"tanggal"`
	NamaObat    string    `json:"nama_obat"`
	HargaSatuan float64   `json:"harga_satuan"`
	Jumlah      int       `json:"jumlah"`
	Subtotal    float64   `json:"subtotal"`
}

type ReportFilterParams struct {
	StartDate string `form:"start_date"` // Format: YYYY-MM-DD
	EndDate   string `form:"end_date"`   // Format: YYYY-MM-DD
	GroupBy   string `form:"group_by"`   // daily, weekly, monthly
}

// US-15: Sales Summary
type SalesSummaryResponse struct {
	TotalPenjualan float64          `json:"total_penjualan"`
	TotalTransaksi int64            `json:"total_transaksi"`
	ChartData      []SalesChartData `json:"chart_data"`
}

type SalesChartData struct {
	Periode         string  `json:"periode"` // Tanggal / Bulan
	TotalPenjualan  float64 `json:"total_penjualan"`
	JumlahTransaksi int64   `json:"jumlah_transaksi"`
}

// US-16: Top & Bottom Medicines
type DrugRankingResponse struct {
	TopMedicines    []DrugSalesStat `json:"top_medicines"`
	BottomMedicines []DrugSalesStat `json:"bottom_medicines"`
}

type DrugSalesStat struct {
	IDObat       uint    `json:"id_obat"`
	NamaObat     string  `json:"nama_obat"`
	TotalTerjual int     `json:"total_terjual"`
	TotalOmset   float64 `json:"total_omset"`
}
