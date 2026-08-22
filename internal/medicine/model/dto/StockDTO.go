package dto

import "time"

type LowStockResponse struct {
	IDObat      uint   `json:"id_obat"`
	NamaObat    string `json:"nama_obat"`
	MerkObat    string `json:"merk_obat"`
	Stok        int    `json:"stok"`
	StokMinimum int    `json:"stok_minimum"`
	NamaJenis   string `json:"nama_jenis"`
}

type ExpiringDrugResponse struct {
	IDObat        uint      `json:"id_obat"`
	NamaObat      string    `json:"nama_obat"`
	MerkObat      string    `json:"merk_obat"`
	TglKadaluarsa time.Time `json:"tgl_kadaluarsa"`
	SisaHari      int       `json:"sisa_hari"`
	Stok          int       `json:"stok"`
	NamaJenis     string    `json:"nama_jenis"`
}

type StockNotificationSummary struct {
	TotalLowStock int64 `json:"total_low_stock"`
	TotalExpiring int64 `json:"total_expiring"`
}