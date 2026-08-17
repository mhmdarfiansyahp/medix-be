package model

import "time"

type Transaksi struct {
	IDTransaksi  uint              `gorm:"primaryKey;column:id_transaksi" json:"id_transaksi"`
	IDUser       uint              `gorm:"column:id_user;not null" json:"id_user"`
	TglTransaksi time.Time         `gorm:"column:tgl_transaksi;autoCreateTime" json:"tgl_transaksi"`
	TotalHarga   float64           `gorm:"column:total_harga;not null" json:"total_harga"`
	Status       int               `gorm:"column:status;default:1;not null" json:"status"`
	CreatedAt    time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Details      []DetailPembelian `gorm:"foreignKey:IDTransaksi;references:IDTransaksi" json:"details"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

type DetailPembelian struct {
	IDDetail    uint    `gorm:"primaryKey;column:id_detail" json:"id_detail"`
	IDTransaksi uint    `gorm:"column:id_transaksi;not null" json:"id_transaksi"`
	IDObat      uint    `gorm:"column:id_obat;not null" json:"id_obat"`
	Jumlah      int     `gorm:"column:jumlah;not null" json:"jumlah"`
	HargaSatuan float64 `gorm:"column:harga_satuan;not null" json:"harga_satuan"`
	Subtotal    float64 `gorm:"column:subtotal;not null" json:"subtotal"`
}

func (DetailPembelian) TableName() string {
	return "detail_pembelian"
}

const (
	StatusTransaksiSelesai   = 1
	StatusTransaksiDibatalkan = 0
)